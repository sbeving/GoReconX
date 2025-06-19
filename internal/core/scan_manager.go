package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorconx/pkg/utils"
)

// ScanManager manages scan execution and provides real-time updates
type ScanManager struct {
	app   *Application
	scans map[string]*ScanExecution
	mutex sync.RWMutex
}

// ScanExecution represents an executing scan
type ScanExecution struct {
	ID          string                 `json:"id"`
	SessionID   string                 `json:"session_id"`
	ModuleName  string                 `json:"module_name"`
	Target      string                 `json:"target"`
	Status      string                 `json:"status"`
	Progress    float64                `json:"progress"`
	StartedAt   int64                  `json:"started_at"`
	CompletedAt int64                  `json:"completed_at,omitempty"`
	Results     map[string]interface{} `json:"results"`
	Error       string                 `json:"error,omitempty"`
	Options     map[string]interface{} `json:"options"`
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewScanManager creates a new scan manager
func NewScanManager(app *Application) *ScanManager {
	return &ScanManager{
		app:   app,
		scans: make(map[string]*ScanExecution),
	}
}

// StartScan starts a new scan execution
func (sm *ScanManager) StartScan(sessionID, moduleName, target string, options map[string]interface{}) (*ScanExecution, error) {
	// Get the module
	module, exists := sm.app.GetModule(moduleName)
	if !exists {
		return nil, fmt.Errorf("module %s not found", moduleName)
	}

	// Create scan execution
	ctx, cancel := context.WithCancel(context.Background())
	scan := &ScanExecution{
		ID:         generateScanID(),
		SessionID:  sessionID,
		ModuleName: moduleName,
		Target:     target,
		Status:     "pending",
		Progress:   0.0,
		StartedAt:  getCurrentTimestamp(),
		Results:    make(map[string]interface{}),
		Options:    options,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Store scan
	sm.mutex.Lock()
	sm.scans[scan.ID] = scan
	sm.mutex.Unlock()

	// Store in database
	sm.storeScanInDB(scan)

	// Start execution in goroutine
	go sm.executeScan(scan, module)

	sm.app.logger.Infof("Started scan %s for module %s on target %s", scan.ID, moduleName, target)
	return scan, nil
}

// GetScan returns a scan by ID
func (sm *ScanManager) GetScan(scanID string) (*ScanExecution, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	scan, exists := sm.scans[scanID]
	return scan, exists
}

// GetSessionScans returns all scans for a session
func (sm *ScanManager) GetSessionScans(sessionID string) []*ScanExecution {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var scans []*ScanExecution
	for _, scan := range sm.scans {
		if scan.SessionID == sessionID {
			scans = append(scans, scan)
		}
	}
	return scans
}

// CancelScan cancels a running scan
func (sm *ScanManager) CancelScan(scanID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	scan, exists := sm.scans[scanID]
	if !exists {
		return fmt.Errorf("scan %s not found", scanID)
	}

	if scan.Status == "running" {
		scan.cancel()
		scan.Status = "cancelled"
		sm.updateScanInDB(scan)
		sm.app.logger.Infof("Cancelled scan %s", scanID)
	}

	return nil
}

// executeScan executes a scan and provides real-time updates
func (sm *ScanManager) executeScan(scan *ScanExecution, module Module) {
	// Update status to running
	scan.Status = "running"
	scan.Progress = 0.1
	sm.updateScanInDB(scan)
	sm.broadcastScanUpdate(scan)

	// Simulate progress updates
	progressTicker := time.NewTicker(1 * time.Second)
	defer progressTicker.Stop()

	// Start progress updates in a separate goroutine
	go func() {
		for {
			select {
			case <-scan.ctx.Done():
				return
			case <-progressTicker.C:
				if scan.Status == "running" && scan.Progress < 0.9 {
					scan.Progress += 0.1
					sm.broadcastScanUpdate(scan)
				}
			}
		}
	}()

	// Execute the actual module
	result, err := module.Execute(scan.Target)

	// Stop progress updates
	scan.cancel()

	// Update final status
	scan.CompletedAt = getCurrentTimestamp()
	if err != nil {
		scan.Status = "failed"
		scan.Error = err.Error()
		sm.app.logger.Errorf("Scan %s failed: %v", scan.ID, err)
	} else {
		scan.Status = "completed"
		scan.Results = map[string]interface{}{
			"data": result,
		}
		sm.app.logger.Infof("Scan %s completed successfully", scan.ID)
	}

	scan.Progress = 1.0
	sm.updateScanInDB(scan)
	sm.broadcastScanUpdate(scan)

	// Store results in session
	session, exists := sm.app.GetSession(scan.SessionID)
	if exists {
		if session.Results == nil {
			session.Results = make(map[string]interface{})
		}
		session.Results[scan.ModuleName] = scan.Results
		sm.app.UpdateSession(session)
	}
}

// storeScanInDB stores scan information in database
func (sm *ScanManager) storeScanInDB(scan *ScanExecution) {
	query := `
		INSERT INTO scans (id, session_id, module_name, target, status, progress, started_at, options)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	optionsJSON := sm.serializeOptions(scan.Options)
	_, err := sm.app.db.Exec(query, scan.ID, scan.SessionID, scan.ModuleName, scan.Target,
		scan.Status, scan.Progress, scan.StartedAt, optionsJSON)
	if err != nil {
		sm.app.logger.Errorf("Failed to store scan in database: %v", err)
	}
}

// updateScanInDB updates scan information in database
func (sm *ScanManager) updateScanInDB(scan *ScanExecution) {
	query := `
		UPDATE scans 
		SET status = ?, progress = ?, completed_at = ?, error_message = ?
		WHERE id = ?
	`
	_, err := sm.app.db.Exec(query, scan.Status, scan.Progress, scan.CompletedAt, scan.Error, scan.ID)
	if err != nil {
		sm.app.logger.Errorf("Failed to update scan in database: %v", err)
	}
}

// broadcastScanUpdate broadcasts scan updates to WebSocket clients
func (sm *ScanManager) broadcastScanUpdate(scan *ScanExecution) {
	// This would integrate with the WebSocket manager
	// For now, we'll just log the update
	sm.app.logger.Infof("Scan %s: %s (%.1f%%)", scan.ID, scan.Status, scan.Progress*100)
}

// serializeOptions serializes options to JSON string
func (sm *ScanManager) serializeOptions(options map[string]interface{}) string {
	// Simple implementation - in production, use proper JSON marshaling
	return "{}"
}

// generateScanID generates a unique scan ID
func generateScanID() string {
	return "scan_" + utils.GenerateRandomString(16)
}

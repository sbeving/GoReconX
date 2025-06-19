package gui

import (
	"time"

	"gorconx/internal/core"
)

// getIndexHTML returns the HTML content for the index page
func getIndexHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoReconX - OSINT & Reconnaissance Platform</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #1a1a2e, #16213e, #0f0f23);
            color: #ffffff;
            min-height: 100vh;
            overflow-x: hidden;
        }
        
        .hero {
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            text-align: center;
            position: relative;
        }
        
        .hero::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: radial-gradient(circle at 20% 80%, rgba(0, 255, 255, 0.1) 0%, transparent 50%),
                        radial-gradient(circle at 80% 20%, rgba(255, 0, 255, 0.1) 0%, transparent 50%);
            z-index: -1;
        }
        
        .logo {
            font-size: 4rem;
            font-weight: bold;
            background: linear-gradient(45deg, #00ffff, #ff00ff, #ffff00, #00ff00);
            background-size: 400% 400%;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            animation: gradientShift 3s ease-in-out infinite;
            margin-bottom: 1rem;
        }
        
        @keyframes gradientShift {
            0%, 100% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
        }
        
        .tagline {
            font-size: 1.5rem;
            margin-bottom: 2rem;
            opacity: 0.9;
            color: #b0b0b0;
        }
        
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 1.5rem;
            margin: 3rem 0;
            max-width: 800px;
        }
        
        .feature {
            background: rgba(255, 255, 255, 0.05);
            padding: 1.5rem;
            border-radius: 10px;
            border: 1px solid rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }
        
        .feature:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 30px rgba(0, 255, 255, 0.2);
        }
        
        .feature h3 {
            color: #00ffff;
            margin-bottom: 0.5rem;
        }
        
        .cta-buttons {
            display: flex;
            gap: 1rem;
            margin-top: 2rem;
            flex-wrap: wrap;
            justify-content: center;
        }
        
        .btn {
            padding: 12px 30px;
            border: none;
            border-radius: 25px;
            font-size: 1rem;
            cursor: pointer;
            transition: all 0.3s ease;
            text-decoration: none;
            display: inline-block;
            font-weight: bold;
        }
        
        .btn-primary {
            background: linear-gradient(45deg, #00ffff, #0080ff);
            color: #000;
        }
        
        .btn-primary:hover {
            transform: scale(1.05);
            box-shadow: 0 5px 20px rgba(0, 255, 255, 0.4);
        }
        
        .btn-secondary {
            background: transparent;
            color: #00ffff;
            border: 2px solid #00ffff;
        }
        
        .btn-secondary:hover {
            background: #00ffff;
            color: #000;
        }
        
        .disclaimer {
            position: fixed;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            background: rgba(255, 0, 0, 0.2);
            color: #ff6b6b;
            padding: 10px 20px;
            border-radius: 20px;
            border: 1px solid rgba(255, 0, 0, 0.3);
            font-size: 0.9rem;
            max-width: 90%;
            text-align: center;
        }
        
        @media (max-width: 768px) {
            .logo { font-size: 2.5rem; }
            .tagline { font-size: 1.2rem; }
            .features { grid-template-columns: 1fr; }
            .cta-buttons { flex-direction: column; align-items: center; }
        }
    </style>
</head>
<body>
    <div class="hero">
        <div class="logo">GoReconX</div>
        <div class="tagline">Advanced OSINT & Reconnaissance Platform</div>
        
        <div class="features">
            <div class="feature">
                <h3>🎯 Modular Design</h3>
                <p>Plug-and-play reconnaissance modules for maximum flexibility</p>
            </div>
            <div class="feature">
                <h3>⚡ High Performance</h3>
                <p>Built with Go for exceptional speed and concurrency</p>
            </div>
            <div class="feature">
                <h3>🔒 Security First</h3>
                <p>Encrypted data storage and secure API management</p>
            </div>
            <div class="feature">
                <h3>📊 Real-time Results</h3>
                <p>Live updates and comprehensive reporting capabilities</p>
            </div>
        </div>
        
        <div class="cta-buttons">
            <a href="/dashboard" class="btn btn-primary">Launch Dashboard</a>
            <a href="/modules" class="btn btn-secondary">Browse Modules</a>
        </div>
    </div>
    
    <div class="disclaimer">
        ⚖️ ETHICAL USE ONLY - Ensure you have explicit permission before scanning any target
    </div>
    
    <script>
        // Add some interactive effects
        document.addEventListener('mousemove', (e) => {
            const cursor = document.createElement('div');
            cursor.style.position = 'fixed';
            cursor.style.left = e.clientX + 'px';
            cursor.style.top = e.clientY + 'px';
            cursor.style.width = '3px';
            cursor.style.height = '3px';
            cursor.style.background = 'rgba(0, 255, 255, 0.5)';
            cursor.style.borderRadius = '50%';
            cursor.style.pointerEvents = 'none';
            cursor.style.zIndex = '9999';
            document.body.appendChild(cursor);
            
            setTimeout(() => {
                cursor.remove();
            }, 100);
        });
    </script>
</body>
</html>`
}

// getDashboardHTML returns the HTML content for the dashboard page
func getDashboardHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard - GoReconX</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #1a1a2e, #16213e);
            color: #ffffff;
            min-height: 100vh;
        }
        
        .navbar {
            background: rgba(0, 0, 0, 0.3);
            backdrop-filter: blur(10px);
            padding: 1rem 2rem;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .navbar-brand {
            font-size: 1.5rem;
            font-weight: bold;
            color: #00ffff;
        }
        
        .navbar-menu {
            display: flex;
            gap: 1rem;
        }
        
        .nav-link {
            color: #fff;
            text-decoration: none;
            padding: 0.5rem 1rem;
            border-radius: 5px;
            transition: background 0.3s;
        }
        
        .nav-link:hover {
            background: rgba(0, 255, 255, 0.1);
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 2rem;
        }
        
        .dashboard-header {
            text-align: center;
            margin-bottom: 3rem;
        }
        
        .dashboard-title {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            background: linear-gradient(45deg, #00ffff, #ff00ff);
            background-size: 400% 400%;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            animation: gradientShift 3s ease-in-out infinite;
        }
        
        @keyframes gradientShift {
            0%, 100% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
        }
        
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 2rem;
            margin-top: 2rem;
        }
        
        .card {
            background: rgba(255, 255, 255, 0.05);
            border-radius: 15px;
            padding: 2rem;
            border: 1px solid rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }
        
        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 15px 35px rgba(0, 255, 255, 0.1);
        }
        
        .card-header {
            display: flex;
            align-items: center;
            margin-bottom: 1rem;
        }
        
        .card-icon {
            font-size: 2rem;
            margin-right: 1rem;
        }
        
        .card h3 {
            color: #00ffff;
            margin-bottom: 1rem;
        }
        
        .btn {
            background: linear-gradient(45deg, #00ffff, #0080ff);
            color: #000;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-weight: bold;
            margin: 5px;
            text-decoration: none;
            display: inline-block;
            transition: all 0.3s ease;
        }
        
        .btn:hover {
            transform: scale(1.05);
            box-shadow: 0 5px 20px rgba(0, 255, 255, 0.4);
        }
        
        .btn-secondary {
            background: transparent;
            color: #00ffff;
            border: 2px solid #00ffff;
        }
        
        .btn-secondary:hover {
            background: #00ffff;
            color: #000;
        }
        
        .btn-danger {
            background: linear-gradient(45deg, #ff6b6b, #ff3333);
            color: #fff;
        }
        
        .status {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.8rem;
            font-weight: bold;
        }
        
        .status.active {
            background: rgba(0, 255, 0, 0.2);
            color: #00ff00;
        }
        
        .status.inactive {
            background: rgba(255, 0, 0, 0.2);
            color: #ff6b6b;
        }
        
        .status.warning {
            background: rgba(255, 255, 0, 0.2);
            color: #ffff00;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin: 1rem 0;
        }
        
        .stat-item {
            text-align: center;
            padding: 1rem;
            background: rgba(0, 0, 0, 0.2);
            border-radius: 10px;
        }
        
        .stat-number {
            font-size: 2rem;
            font-weight: bold;
            color: #00ffff;
        }
        
        .stat-label {
            font-size: 0.9rem;
            opacity: 0.8;
        }
        
        .quick-actions {
            display: flex;
            flex-wrap: wrap;
            gap: 1rem;
            margin-top: 1rem;
        }
        
        .module-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
            gap: 1rem;
            margin-top: 1rem;
        }
        
        .module-card {
            background: rgba(0, 0, 0, 0.2);
            padding: 1rem;
            border-radius: 10px;
            border: 1px solid rgba(255, 255, 255, 0.1);
            transition: all 0.3s ease;
        }
        
        .module-card:hover {
            background: rgba(0, 255, 255, 0.1);
            border-color: #00ffff;
        }
        
        .module-name {
            font-weight: bold;
            color: #00ffff;
            margin-bottom: 0.5rem;
        }
        
        .module-description {
            font-size: 0.9rem;
            opacity: 0.8;
            margin-bottom: 1rem;
        }
        
        .module-tags {
            display: flex;
            flex-wrap: wrap;
            gap: 0.25rem;
        }
        
        .tag {
            background: rgba(0, 255, 255, 0.2);
            color: #00ffff;
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 0.7rem;
        }
        
        .real-time-feed {
            max-height: 300px;
            overflow-y: auto;
            background: rgba(0, 0, 0, 0.3);
            border-radius: 10px;
            padding: 1rem;
        }
        
        .feed-item {
            padding: 0.5rem;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            font-size: 0.9rem;
        }
        
        .feed-time {
            color: #888;
            font-size: 0.8rem;
        }
        
        .loading {
            text-align: center;
            color: #888;
            font-style: italic;
        }
        
        @media (max-width: 768px) {
            .container { padding: 1rem; }
            .grid { grid-template-columns: 1fr; }
            .stats-grid { grid-template-columns: repeat(2, 1fr); }
            .navbar { flex-direction: column; gap: 1rem; }
            .dashboard-title { font-size: 2rem; }
        }
    </style>
</head>
<body>
    <nav class="navbar">
        <div class="navbar-brand">GoReconX Dashboard</div>
        <div class="navbar-menu">
            <a href="/" class="nav-link">Home</a>
            <a href="/dashboard" class="nav-link">Dashboard</a>
            <a href="/modules" class="nav-link">Modules</a>
            <a href="/sessions" class="nav-link">Sessions</a>
            <a href="/reports" class="nav-link">Reports</a>
            <a href="/settings" class="nav-link">Settings</a>
        </div>
    </nav>
    
    <div class="container">
        <div class="dashboard-header">
            <h1 class="dashboard-title">OSINT & Reconnaissance Dashboard</h1>
            <p>Comprehensive intelligence gathering and network reconnaissance platform</p>
        </div>
        
        <div class="grid">
            <!-- Quick Start Card -->
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">🚀</div>
                    <h3>Quick Start</h3>
                </div>
                <p>Launch a new reconnaissance session or continue existing work</p>
                <div class="quick-actions">
                    <button class="btn" onclick="showNewSessionModal()">New Session</button>
                    <a href="/sessions" class="btn btn-secondary">View Sessions</a>
                </div>
            </div>
            
            <!-- System Status Card -->
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">⚙️</div>
                    <h3>System Status</h3>
                </div>
                <div class="stats-grid">
                    <div class="stat-item">
                        <div class="stat-number" id="module-count">0</div>
                        <div class="stat-label">Modules</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-number" id="session-count">0</div>
                        <div class="stat-label">Sessions</div>
                    </div>
                </div>
                <p>GUI Server: <span class="status active">Active</span></p>
                <p>Database: <span class="status active">Connected</span></p>
                <p>API Server: <span class="status active">Running</span></p>
            </div>
            
            <!-- Modules Overview -->
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">🔧</div>
                    <h3>Available Modules</h3>
                </div>
                <div id="modules-overview" class="loading">Loading modules...</div>
                <a href="/modules" class="btn">Browse All Modules</a>
            </div>
            
            <!-- Recent Activity -->
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">📊</div>
                    <h3>Recent Activity</h3>
                </div>
                <div id="recent-activity" class="real-time-feed">
                    <div class="loading">Loading recent activity...</div>
                </div>
            </div>
            
            <!-- Quick Tools -->
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">🛠️</div>
                    <h3>Quick Tools</h3>
                </div>
                <p>Access commonly used reconnaissance tools</p>
                <div class="quick-actions">
                    <button class="btn" onclick="quickDomainScan()">Domain Scan</button>
                    <button class="btn" onclick="quickPortScan()">Port Scan</button>
                    <button class="btn" onclick="quickWebScan()">Web Scan</button>
                </div>
            </div>
            
            <!-- Security Notice -->
            <div class="card">
                <div class="card-header">
                    <div class="card-icon">🔒</div>
                    <h3>Security & Ethics</h3>
                </div>
                <p style="color: #ff6b6b; font-weight: bold;">⚖️ ETHICAL USE ONLY</p>
                <p>Always ensure you have explicit permission before scanning any target. This tool is for legitimate security assessments, educational purposes, and authorized penetration testing only.</p>
                <button class="btn btn-secondary" onclick="showEthicsGuidelines()">Ethics Guidelines</button>
            </div>
        </div>
    </div>
    
    <!-- New Session Modal -->
    <div id="newSessionModal" style="display: none; position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.8); z-index: 1000;">
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); background: linear-gradient(135deg, #1a1a2e, #16213e); padding: 2rem; border-radius: 15px; border: 1px solid rgba(255, 255, 255, 0.1); min-width: 400px;">
            <h3 style="color: #00ffff; margin-bottom: 1rem;">Create New Session</h3>
            <form id="newSessionForm">
                <div style="margin-bottom: 1rem;">
                    <label style="display: block; margin-bottom: 0.5rem;">Session Name:</label>
                    <input type="text" id="sessionName" style="width: 100%; padding: 0.5rem; border: 1px solid rgba(255,255,255,0.3); background: rgba(255,255,255,0.1); color: #fff; border-radius: 5px;" required>
                </div>
                <div style="margin-bottom: 1rem;">
                    <label style="display: block; margin-bottom: 0.5rem;">Target (Domain/IP):</label>
                    <input type="text" id="sessionTarget" style="width: 100%; padding: 0.5rem; border: 1px solid rgba(255,255,255,0.3); background: rgba(255,255,255,0.1); color: #fff; border-radius: 5px;" required>
                </div>
                <div style="display: flex; gap: 1rem; justify-content: flex-end;">
                    <button type="button" class="btn btn-secondary" onclick="closeNewSessionModal()">Cancel</button>
                    <button type="submit" class="btn">Create Session</button>
                </div>
            </form>
        </div>
    </div>
    
    <script>
        let ws = null;
        let reconnectTimer = null;
        
        // Initialize dashboard
        document.addEventListener('DOMContentLoaded', function() {
            loadModules();
            loadSessions();
            loadRecentActivity();
            connectWebSocket();
        });
        
        // WebSocket connection for real-time updates
        function connectWebSocket() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws';
            
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                console.log('WebSocket connected');
                if (reconnectTimer) {
                    clearInterval(reconnectTimer);
                    reconnectTimer = null;
                }
            };
            
            ws.onmessage = function(event) {
                const data = JSON.parse(event.data);
                handleRealTimeUpdate(data);
            };
            
            ws.onclose = function() {
                console.log('WebSocket disconnected');
                // Attempt to reconnect every 5 seconds
                if (!reconnectTimer) {
                    reconnectTimer = setInterval(connectWebSocket, 5000);
                }
            };
            
            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
            };
        }
        
        function handleRealTimeUpdate(data) {
            // Update recent activity feed
            const feed = document.getElementById('recent-activity');
            const item = document.createElement('div');
            item.className = 'feed-item';
            item.innerHTML = '<div class="feed-time">' + new Date().toLocaleTimeString() + '</div>' +
                           '<div>' + data.message + '</div>';
            feed.insertBefore(item, feed.firstChild);
            
            // Keep only last 10 items
            while (feed.children.length > 10) {
                feed.removeChild(feed.lastChild);
            }
        }
        
        async function loadModules() {
            try {
                const response = await fetch('/api/modules');
                const modules = await response.json();
                
                document.getElementById('module-count').textContent = Object.keys(modules).length;
                
                const overview = document.getElementById('modules-overview');
                overview.innerHTML = '';
                
                const moduleGrid = document.createElement('div');
                moduleGrid.className = 'module-grid';
                
                for (const [name, info] of Object.entries(modules)) {
                    const moduleCard = document.createElement('div');
                    moduleCard.className = 'module-card';
                    moduleCard.innerHTML = 
                        '<div class="module-name">' + name + '</div>' +
                        '<div class="module-description">' + (info.description || 'No description available') + '</div>' +
                        '<div class="module-tags">' +
                        (info.tags || []).map(tag => '<span class="tag">' + tag + '</span>').join('') +
                        '</div>';
                    moduleGrid.appendChild(moduleCard);
                }
                
                overview.appendChild(moduleGrid);
            } catch (error) {
                console.error('Failed to load modules:', error);
                document.getElementById('modules-overview').innerHTML = '<div style="color: #ff6b6b;">Failed to load modules</div>';
            }
        }
        
        async function loadSessions() {
            try {
                const response = await fetch('/api/sessions');
                const sessions = await response.json();
                document.getElementById('session-count').textContent = sessions.length;
            } catch (error) {
                console.error('Failed to load sessions:', error);
            }
        }
        
        async function loadRecentActivity() {
            // Simulate loading recent activity
            const activities = [
                'System started successfully',
                'Modules loaded: 5 modules available',
                'Database connection established',
                'Ready for reconnaissance operations'
            ];
            
            const feed = document.getElementById('recent-activity');
            feed.innerHTML = '';
            
            activities.forEach((activity, index) => {
                setTimeout(() => {
                    const item = document.createElement('div');
                    item.className = 'feed-item';
                    const time = new Date(Date.now() - (activities.length - index) * 1000);
                    item.innerHTML = '<div class="feed-time">' + time.toLocaleTimeString() + '</div>' +
                                   '<div>' + activity + '</div>';
                    feed.appendChild(item);
                }, index * 200);
            });
        }
        
        function showNewSessionModal() {
            document.getElementById('newSessionModal').style.display = 'block';
            document.getElementById('sessionName').focus();
        }
        
        function closeNewSessionModal() {
            document.getElementById('newSessionModal').style.display = 'none';
            document.getElementById('newSessionForm').reset();
        }
        
        document.getElementById('newSessionForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const name = document.getElementById('sessionName').value;
            const target = document.getElementById('sessionTarget').value;
            
            try {
                const response = await fetch('/api/sessions', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name, target })
                });
                
                if (response.ok) {
                    const session = await response.json();
                    alert('Session "' + session.name + '" created successfully!');
                    closeNewSessionModal();
                    loadSessions();
                    
                    // Redirect to session page
                    window.location.href = '/sessions/' + session.id;
                } else {
                    throw new Error('Failed to create session');
                }
            } catch (error) {
                alert('Failed to create session: ' + error.message);
            }
        });
        
        function quickDomainScan() {
            const domain = prompt('Enter domain to scan:');
            if (domain) {
                // Create quick session and redirect
                fetch('/api/sessions', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 
                        name: 'Quick Domain Scan - ' + domain, 
                        target: domain 
                    })
                }).then(response => response.json())
                .then(session => {
                    window.location.href = '/sessions/' + session.id + '?module=domain_enum';
                });
            }
        }
        
        function quickPortScan() {
            const target = prompt('Enter IP/hostname to scan:');
            if (target) {
                fetch('/api/sessions', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 
                        name: 'Quick Port Scan - ' + target, 
                        target: target 
                    })
                }).then(response => response.json())
                .then(session => {
                    window.location.href = '/sessions/' + session.id + '?module=port_scan';
                });
            }
        }
        
        function quickWebScan() {
            const url = prompt('Enter URL to scan:');
            if (url) {
                fetch('/api/sessions', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 
                        name: 'Quick Web Scan - ' + url, 
                        target: url 
                    })
                }).then(response => response.json())
                .then(session => {
                    window.location.href = '/sessions/' + session.id + '?module=web_enum';
                });
            }
        }
        
        function showEthicsGuidelines() {
            alert('ETHICAL USE GUIDELINES:\\n\\n' +
                  '1. Always obtain explicit written permission before scanning any target\\n' +
                  '2. Only use on systems you own or have been authorized to test\\n' +
                  '3. Respect rate limits and avoid causing service disruption\\n' +
                  '4. Follow all applicable laws and regulations\\n' +
                  '5. Report findings responsibly through proper channels\\n' +
                  '6. Use for legitimate security research and education only\\n\\n' +
                  'Unauthorized scanning is illegal and unethical.');
        }
        
        // Close modal when clicking outside
        document.getElementById('newSessionModal').addEventListener('click', function(e) {
            if (e.target === this) {
                closeNewSessionModal();
            }
        });
        
        // Keyboard shortcuts
        document.addEventListener('keydown', function(e) {
            if (e.ctrlKey && e.key === 'n') {
                e.preventDefault();
                showNewSessionModal();
            }
            if (e.key === 'Escape') {
                closeNewSessionModal();
            }
        });
    </script>
</body>
</html>`
}

// getModulesHTML returns the HTML content for the modules page
func getModulesHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Modules - GoReconX</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #1a1a2e, #16213e);
            color: #ffffff;
            min-height: 100vh;
        }
        
        .navbar {
            background: rgba(0, 0, 0, 0.3);
            backdrop-filter: blur(10px);
            padding: 1rem 2rem;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .navbar-brand {
            font-size: 1.5rem;
            font-weight: bold;
            color: #00ffff;
        }
        
        .navbar-menu {
            display: flex;
            gap: 1rem;
        }
        
        .nav-link {
            color: #fff;
            text-decoration: none;
            padding: 0.5rem 1rem;
            border-radius: 5px;
            transition: background 0.3s;
        }
        
        .nav-link:hover {
            background: rgba(0, 255, 255, 0.1);
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 2rem;
        }
        
        .page-header {
            text-align: center;
            margin-bottom: 3rem;
        }
        
        .page-title {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            background: linear-gradient(45deg, #00ffff, #ff00ff);
            background-size: 400% 400%;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            animation: gradientShift 3s ease-in-out infinite;
        }
        
        @keyframes gradientShift {
            0%, 100% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
        }
        
        .filters {
            display: flex;
            flex-wrap: wrap;
            gap: 1rem;
            margin-bottom: 2rem;
            align-items: center;
        }
        
        .filter-group {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        
        .filter-label {
            font-weight: bold;
            color: #00ffff;
        }
        
        .filter-select {
            background: rgba(255, 255, 255, 0.1);
            border: 1px solid rgba(255, 255, 255, 0.3);
            color: #fff;
            padding: 0.5rem;
            border-radius: 5px;
        }
        
        .search-box {
            background: rgba(255, 255, 255, 0.1);
            border: 1px solid rgba(255, 255, 255, 0.3);
            color: #fff;
            padding: 0.5rem;
            border-radius: 5px;
            width: 300px;
        }
        
        .modules-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
            gap: 2rem;
        }
        
        .module-card {
            background: rgba(255, 255, 255, 0.05);
            border-radius: 15px;
            padding: 2rem;
            border: 1px solid rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }
        
        .module-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 15px 35px rgba(0, 255, 255, 0.1);
        }
        
        .module-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 1rem;
        }
        
        .module-icon {
            font-size: 2.5rem;
            margin-bottom: 1rem;
        }
        
        .module-name {
            font-size: 1.5rem;
            font-weight: bold;
            color: #00ffff;
            margin-bottom: 0.5rem;
        }
        
        .module-category {
            background: rgba(0, 255, 255, 0.2);
            color: #00ffff;
            padding: 0.25rem 0.75rem;
            border-radius: 12px;
            font-size: 0.8rem;
            font-weight: bold;
        }
        
        .module-description {
            margin-bottom: 1rem;
            opacity: 0.9;
            line-height: 1.5;
        }
        
        .module-tags {
            display: flex;
            flex-wrap: wrap;
            gap: 0.5rem;
            margin-bottom: 1.5rem;
        }
        
        .tag {
            background: rgba(255, 255, 255, 0.1);
            color: #fff;
            padding: 0.25rem 0.5rem;
            border-radius: 8px;
            font-size: 0.8rem;
        }
        
        .module-options {
            margin-bottom: 1.5rem;
        }
        
        .options-title {
            font-weight: bold;
            color: #00ffff;
            margin-bottom: 0.5rem;
        }
        
        .option-item {
            display: flex;
            justify-content: space-between;
            padding: 0.25rem 0;
            font-size: 0.9rem;
        }
        
        .option-name {
            font-weight: bold;
        }
        
        .option-type {
            color: #888;
        }
        
        .module-actions {
            display: flex;
            gap: 1rem;
        }
        
        .btn {
            background: linear-gradient(45deg, #00ffff, #0080ff);
            color: #000;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-weight: bold;
            text-decoration: none;
            display: inline-block;
            transition: all 0.3s ease;
            text-align: center;
        }
        
        .btn:hover {
            transform: scale(1.05);
            box-shadow: 0 5px 20px rgba(0, 255, 255, 0.4);
        }
        
        .btn-secondary {
            background: transparent;
            color: #00ffff;
            border: 2px solid #00ffff;
        }
        
        .btn-secondary:hover {
            background: #00ffff;
            color: #000;
        }
        
        .module-stats {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 1rem;
            padding-top: 1rem;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
            font-size: 0.9rem;
            opacity: 0.8;
        }
        
        .category-passive_osint .module-icon { color: #4CAF50; }
        .category-active_recon .module-icon { color: #FF9800; }
        .category-network .module-icon { color: #2196F3; }
        .category-web .module-icon { color: #9C27B0; }
        .category-social .module-icon { color: #E91E63; }
        
        .loading {
            text-align: center;
            color: #888;
            font-style: italic;
            padding: 2rem;
        }
        
        .no-modules {
            text-align: center;
            color: #888;
            padding: 2rem;
        }
        
        @media (max-width: 768px) {
            .container { padding: 1rem; }
            .modules-grid { grid-template-columns: 1fr; }
            .filters { flex-direction: column; align-items: stretch; }
            .search-box { width: 100%; }
            .page-title { font-size: 2rem; }
        }
    </style>
</head>
<body>
    <nav class="navbar">
        <div class="navbar-brand">GoReconX Modules</div>
        <div class="navbar-menu">
            <a href="/" class="nav-link">Home</a>
            <a href="/dashboard" class="nav-link">Dashboard</a>
            <a href="/modules" class="nav-link">Modules</a>
            <a href="/sessions" class="nav-link">Sessions</a>
            <a href="/reports" class="nav-link">Reports</a>
            <a href="/settings" class="nav-link">Settings</a>
        </div>
    </nav>
    
    <div class="container">
        <div class="page-header">
            <h1 class="page-title">Reconnaissance Modules</h1>
            <p>Comprehensive collection of OSINT and reconnaissance tools</p>
        </div>
        
        <div class="filters">
            <div class="filter-group">
                <span class="filter-label">Category:</span>
                <select id="categoryFilter" class="filter-select">
                    <option value="">All Categories</option>
                    <option value="passive_osint">Passive OSINT</option>
                    <option value="active_recon">Active Reconnaissance</option>
                    <option value="network">Network Analysis</option>
                    <option value="web">Web Enumeration</option>
                    <option value="social">Social Engineering</option>
                </select>
            </div>
            
            <div class="filter-group">
                <span class="filter-label">Search:</span>
                <input type="text" id="searchBox" class="search-box" placeholder="Search modules...">
            </div>
        </div>
        
        <div id="modulesContainer" class="modules-grid">
            <div class="loading">Loading modules...</div>
        </div>
    </div>
    
    <script>
        let allModules = {};
        
        document.addEventListener('DOMContentLoaded', function() {
            loadModules();
            
            // Set up filters
            document.getElementById('categoryFilter').addEventListener('change', filterModules);
            document.getElementById('searchBox').addEventListener('input', filterModules);
        });
        
        async function loadModules() {
            try {
                const response = await fetch('/api/modules');
                allModules = await response.json();
                displayModules(allModules);
            } catch (error) {
                console.error('Failed to load modules:', error);
                document.getElementById('modulesContainer').innerHTML = 
                    '<div class="no-modules">Failed to load modules. Please try again.</div>';
            }
        }
        
        function displayModules(modules) {
            const container = document.getElementById('modulesContainer');
            
            if (Object.keys(modules).length === 0) {
                container.innerHTML = '<div class="no-modules">No modules found.</div>';
                return;
            }
            
            container.innerHTML = '';
            
            for (const [name, module] of Object.entries(modules)) {
                const moduleCard = createModuleCard(name, module);
                container.appendChild(moduleCard);
            }
        }
        
        function createModuleCard(name, module) {
            const card = document.createElement('div');
            card.className = 'module-card category-' + (module.category || 'unknown');
            
            const icon = getModuleIcon(module.category);
            const tags = (module.tags || []).map(tag => '<span class="tag">' + tag + '</span>').join('');
            
            let optionsHTML = '';
            if (module.options && module.options.length > 0) {
                optionsHTML = '<div class="module-options">' +
                    '<div class="options-title">Configuration Options:</div>';
                module.options.slice(0, 3).forEach(option => {
                    optionsHTML += '<div class="option-item">' +
                        '<span class="option-name">' + option.name + '</span>' +
                        '<span class="option-type">' + option.type + '</span>' +
                        '</div>';
                });
                if (module.options.length > 3) {
                    optionsHTML += '<div style="color: #888; font-size: 0.8rem;">... and ' + 
                        (module.options.length - 3) + ' more options</div>';
                }
                optionsHTML += '</div>';
            }
            
            card.innerHTML = 
                '<div class="module-header">' +
                    '<div class="module-icon">' + icon + '</div>' +
                    '<div class="module-category">' + (module.category || 'unknown').replace('_', ' ') + '</div>' +
                '</div>' +
                '<div class="module-name">' + name + '</div>' +
                '<div class="module-description">' + (module.description || 'No description available') + '</div>' +
                '<div class="module-tags">' + tags + '</div>' +
                optionsHTML +
                '<div class="module-actions">' +
                    '<button class="btn" onclick="runModule(\'' + name + '\')">Run Module</button>' +
                    '<button class="btn btn-secondary" onclick="configureModule(\'' + name + '\')">Configure</button>' +
                '</div>' +
                '<div class="module-stats">' +
                    '<span>Version: ' + (module.version || '1.0.0') + '</span>' +
                    '<span>Author: ' + (module.author || 'Unknown') + '</span>' +
                '</div>';
            
            return card;
        }
        
        function getModuleIcon(category) {
            const icons = {
                'passive_osint': '🔍',
                'active_recon': '🎯',
                'network': '🌐',
                'web': '💻',
                'social': '👥',
                'unknown': '🔧'
            };
            return icons[category] || icons.unknown;
        }
        
        function filterModules() {
            const category = document.getElementById('categoryFilter').value;
            const search = document.getElementById('searchBox').value.toLowerCase();
            
            let filteredModules = {};
            
            for (const [name, module] of Object.entries(allModules)) {
                let matches = true;
                
                // Category filter
                if (category && module.category !== category) {
                    matches = false;
                }
                
                // Search filter
                if (search) {
                    const searchText = (name + ' ' + 
                        (module.description || '') + ' ' + 
                        (module.tags || []).join(' ')).toLowerCase();
                    if (!searchText.includes(search)) {
                        matches = false;
                    }
                }
                
                if (matches) {
                    filteredModules[name] = module;
                }
            }
            
            displayModules(filteredModules);
        }
        
        function runModule(moduleName) {
            const target = prompt('Enter target (domain/IP/URL):');
            if (target) {
                // Create a new session for this module
                fetch('/api/sessions', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 
                        name: moduleName + ' - ' + target, 
                        target: target 
                    })
                }).then(response => response.json())
                .then(session => {
                    window.location.href = '/sessions/' + session.id + '?module=' + moduleName;
                }).catch(error => {
                    alert('Failed to create session: ' + error.message);
                });
            }
        }
        
        function configureModule(moduleName) {
            alert('Module configuration will be available in the next update!');
        }
    </script>
</body>
</html>`
}

// getSessionsHTML returns the HTML content for the sessions page
func getSessionsHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sessions - GoReconX</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #1a1a2e, #16213e);
            color: #ffffff;
            min-height: 100vh;
        }
        
        .navbar {
            background: rgba(0, 0, 0, 0.3);
            backdrop-filter: blur(10px);
            padding: 1rem 2rem;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .navbar-brand {
            font-size: 1.5rem;
            font-weight: bold;
            color: #00ffff;
        }
        
        .navbar-menu {
            display: flex;
            gap: 1rem;
        }
        
        .nav-link {
            color: #fff;
            text-decoration: none;
            padding: 0.5rem 1rem;
            border-radius: 5px;
            transition: background 0.3s;
        }
        
        .nav-link:hover {
            background: rgba(0, 255, 255, 0.1);
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 2rem;
        }
        
        .page-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 2rem;
        }
        
        .page-title {
            font-size: 2.5rem;
            background: linear-gradient(45deg, #00ffff, #ff00ff);
            background-size: 400% 400%;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            animation: gradientShift 3s ease-in-out infinite;
        }
        
        @keyframes gradientShift {
            0%, 100% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
        }
        
        .btn {
            background: linear-gradient(45deg, #00ffff, #0080ff);
            color: #000;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-weight: bold;
            text-decoration: none;
            display: inline-block;
            transition: all 0.3s ease;
        }
        
        .btn:hover {
            transform: scale(1.05);
            box-shadow: 0 5px 20px rgba(0, 255, 255, 0.4);
        }
        
        .sessions-table {
            background: rgba(255, 255, 255, 0.05);
            border-radius: 15px;
            overflow: hidden;
            border: 1px solid rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
        }
        
        .table {
            width: 100%;
            border-collapse: collapse;
        }
        
        .table th {
            background: rgba(0, 255, 255, 0.1);
            padding: 1rem;
            text-align: left;
            font-weight: bold;
            color: #00ffff;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        }
        
        .table td {
            padding: 1rem;
            border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        }
        
        .table tr:hover {
            background: rgba(255, 255, 255, 0.05);
        }
        
        .status {
            display: inline-block;
            padding: 0.25rem 0.75rem;
            border-radius: 12px;
            font-size: 0.8rem;
            font-weight: bold;
        }
        
        .status.created {
            background: rgba(0, 255, 255, 0.2);
            color: #00ffff;
        }
        
        .status.running {
            background: rgba(255, 255, 0, 0.2);
            color: #ffff00;
        }
        
        .status.completed {
            background: rgba(0, 255, 0, 0.2);
            color: #00ff00;
        }
        
        .status.error {
            background: rgba(255, 0, 0, 0.2);
            color: #ff6b6b;
        }
        
        .action-buttons {
            display: flex;
            gap: 0.5rem;
        }
        
        .btn-small {
            padding: 0.5rem 1rem;
            font-size: 0.8rem;
        }
        
        .btn-secondary {
            background: transparent;
            color: #00ffff;
            border: 2px solid #00ffff;
        }
        
        .btn-secondary:hover {
            background: #00ffff;
            color: #000;
        }
        
        .btn-danger {
            background: linear-gradient(45deg, #ff6b6b, #ff3333);
            color: #fff;
        }
        
        .loading {
            text-align: center;
            color: #888;
            font-style: italic;
            padding: 2rem;
        }
        
        .no-sessions {
            text-align: center;
            color: #888;
            padding: 3rem;
        }
        
        .no-sessions h3 {
            margin-bottom: 1rem;
        }
        
        @media (max-width: 768px) {
            .container { padding: 1rem; }
            .page-header { flex-direction: column; gap: 1rem; }
            .table { font-size: 0.9rem; }
            .table th, .table td { padding: 0.5rem; }
        }
    </style>
</head>
<body>
    <nav class="navbar">
        <div class="navbar-brand">GoReconX Sessions</div>
        <div class="navbar-menu">
            <a href="/" class="nav-link">Home</a>
            <a href="/dashboard" class="nav-link">Dashboard</a>
            <a href="/modules" class="nav-link">Modules</a>
            <a href="/sessions" class="nav-link">Sessions</a>
            <a href="/reports" class="nav-link">Reports</a>
            <a href="/settings" class="nav-link">Settings</a>
        </div>
    </nav>
    
    <div class="container">
        <div class="page-header">
            <h1 class="page-title">Reconnaissance Sessions</h1>
            <button class="btn" onclick="showNewSessionModal()">New Session</button>
        </div>
        
        <div class="sessions-table">
            <table class="table">
                <thead>
                    <tr>
                        <th>Session Name</th>
                        <th>Target</th>
                        <th>Status</th>
                        <th>Created</th>
                        <th>Updated</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody id="sessionsTableBody">
                    <tr>
                        <td colspan="6" class="loading">Loading sessions...</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    
    <script>
        document.addEventListener('DOMContentLoaded', function() {
            loadSessions();
            
            // Refresh sessions every 30 seconds
            setInterval(loadSessions, 30000);
        });
        
        async function loadSessions() {
            try {
                const response = await fetch('/api/sessions');
                const sessions = await response.json();
                displaySessions(sessions);
            } catch (error) {
                console.error('Failed to load sessions:', error);
                document.getElementById('sessionsTableBody').innerHTML = 
                    '<tr><td colspan="6" style="text-align: center; color: #ff6b6b;">Failed to load sessions</td></tr>';
            }
        }
        
        function displaySessions(sessions) {
            const tbody = document.getElementById('sessionsTableBody');
            
            if (sessions.length === 0) {
                tbody.innerHTML = 
                    '<tr><td colspan="6" class="no-sessions">' +
                    '<h3>No sessions yet</h3>' +
                    '<p>Create your first reconnaissance session to get started</p>' +
                    '<button class="btn" onclick="showNewSessionModal()">Create Session</button>' +
                    '</td></tr>';
                return;
            }
            
            tbody.innerHTML = '';
            
            sessions.forEach(session => {
                const row = document.createElement('tr');
                row.innerHTML = 
                    '<td><strong>' + session.name + '</strong></td>' +
                    '<td>' + session.target + '</td>' +
                    '<td><span class="status ' + session.status + '">' + session.status + '</span></td>' +
                    '<td>' + formatDate(session.created_at) + '</td>' +
                    '<td>' + formatDate(session.updated_at) + '</td>' +
                    '<td>' +
                        '<div class="action-buttons">' +
                            '<button class="btn btn-small" onclick="openSession(\'' + session.id + '\')">Open</button>' +
                            '<button class="btn btn-secondary btn-small" onclick="duplicateSession(\'' + session.id + '\')">Duplicate</button>' +
                            '<button class="btn btn-danger btn-small" onclick="deleteSession(\'' + session.id + '\')">Delete</button>' +
                        '</div>' +
                    '</td>';
                tbody.appendChild(row);
            });
        }
        
        function formatDate(timestamp) {
            if (!timestamp) return 'N/A';
            const date = new Date(timestamp * 1000);
            return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
        }
        
        function showNewSessionModal() {
            // This would show the same modal as in dashboard
            const name = prompt('Session name:');
            const target = prompt('Target (domain/IP/URL):');
            
            if (name && target) {
                createSession(name, target);
            }
        }
        
        async function createSession(name, target) {
            try {
                const response = await fetch('/api/sessions', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name, target })
                });
                
                if (response.ok) {
                    const session = await response.json();
                    alert('Session "' + session.name + '" created successfully!');
                    loadSessions();
                } else {
                    throw new Error('Failed to create session');
                }
            } catch (error) {
                alert('Failed to create session: ' + error.message);
            }
        }
        
        function openSession(sessionId) {
            window.location.href = '/sessions/' + sessionId;
        }
        
        async function duplicateSession(sessionId) {
            try {
                const response = await fetch('/api/sessions/' + sessionId);
                const session = await response.json();
                
                const newName = prompt('New session name:', session.name + ' (Copy)');
                if (newName) {
                    createSession(newName, session.target);
                }
            } catch (error) {
                alert('Failed to duplicate session: ' + error.message);
            }
        }
        
        async function deleteSession(sessionId) {
            if (confirm('Are you sure you want to delete this session? This action cannot be undone.')) {
                try {
                    const response = await fetch('/api/sessions/' + sessionId, {
                        method: 'DELETE'
                    });
                    
                    if (response.ok) {
                        alert('Session deleted successfully');
                        loadSessions();
                    } else {
                        throw new Error('Failed to delete session');
                    }
                } catch (error) {
                    alert('Failed to delete session: ' + error.message);
                }
            }
        }
    </script>
</body>
</html>`
}

// getSessionDetailHTML returns HTML for individual session detail page
func getSessionDetailHTML(session *core.Session) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Session: ` + session.Name + ` - GoReconX</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            margin: 0;
            padding: 20px;
            min-height: 100vh;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            padding: 30px;
            box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
        }
        .header {
            border-bottom: 2px solid #667eea;
            padding-bottom: 20px;
            margin-bottom: 30px;
        }
        .header h1 {
            color: #2c3e50;
            margin: 0;
            font-size: 2rem;
        }
        .session-info {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .info-card {
            background: #f8f9fa;
            border-radius: 10px;
            padding: 20px;
        }
        .info-card h3 {
            color: #667eea;
            margin: 0 0 10px 0;
        }
        .nav-back {
            background: #667eea;
            color: white;
            padding: 10px 20px;
            border-radius: 10px;
            text-decoration: none;
            display: inline-block;
            margin-bottom: 20px;
        }
        .nav-back:hover {
            background: #5a6fd8;
        }
        .results-section {
            margin-top: 30px;
        }
        .results-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
        }
        .module-result {
            background: #f8f9fa;
            border-radius: 10px;
            padding: 20px;
            border-left: 4px solid #667eea;
        }
        .module-result h4 {
            color: #2c3e50;
            margin: 0 0 15px 0;
        }
        .json-viewer {
            background: #2c3e50;
            color: #ecf0f1;
            padding: 15px;
            border-radius: 8px;
            font-family: 'Courier New', monospace;
            white-space: pre-wrap;
            max-height: 300px;
            overflow-y: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/sessions" class="nav-back">← Back to Sessions</a>
        
        <div class="header">
            <h1>📊 ` + session.Name + `</h1>
        </div>

        <div class="session-info">
            <div class="info-card">
                <h3>🎯 Target</h3>
                <p>` + session.Target + `</p>
            </div>
            <div class="info-card">
                <h3>📅 Created</h3>
                <p><span id="created-date">` + formatTimestamp(session.CreatedAt) + `</span></p>
            </div>
            <div class="info-card">
                <h3>🔄 Status</h3>
                <p><span style="color: #27ae60;">` + session.Status + `</span></p>
            </div>
            <div class="info-card">
                <h3>🆔 Session ID</h3>
                <p><code>` + session.ID + `</code></p>
            </div>
        </div>

        <div class="results-section">
            <h2>📋 Scan Results</h2>
            <div class="results-grid" id="results-container">
                <!-- Results will be loaded here -->
            </div>
        </div>
    </div>

    <script>
        // Load session results
        async function loadResults() {
            try {
                const response = await fetch('/api/scans?session_id=` + session.ID + `');
                const scans = await response.json();
                
                const container = document.getElementById('results-container');
                
                if (scans.length === 0) {
                    container.innerHTML = '<p style="text-align: center; color: #7f8c8d;">No scan results yet.</p>';
                    return;
                }
                
                container.innerHTML = scans.map(scan => ` + "`" + `
                    <div class="module-result">
                        <h4>${scan.module_name}</h4>
                        <p><strong>Status:</strong> <span style="color: ${getStatusColor(scan.status)}">${scan.status}</span></p>
                        <p><strong>Progress:</strong> ${Math.round(scan.progress * 100)}%</p>
                        ${scan.results ? ` + "`" + `<div class="json-viewer">${JSON.stringify(scan.results, null, 2)}</div>` + "`" + ` : ''}
                        ${scan.error ? ` + "`" + `<p style="color: #e74c3c;"><strong>Error:</strong> ${scan.error}</p>` + "`" + ` : ''}
                    </div>
                ` + "`" + `).join('');
            } catch (error) {
                console.error('Failed to load results:', error);
                document.getElementById('results-container').innerHTML = 
                    '<p style="text-align: center; color: #e74c3c;">Failed to load results.</p>';
            }
        }
        
        function getStatusColor(status) {
            switch(status) {
                case 'completed': return '#27ae60';
                case 'running': return '#f39c12';
                case 'failed': return '#e74c3c';
                case 'cancelled': return '#95a5a6';
                default: return '#3498db';
            }
        }
        
        // Load results on page load
        document.addEventListener('DOMContentLoaded', loadResults);
        
        // Refresh every 5 seconds if there are running scans
        setInterval(loadResults, 5000);
    </script>
</body>
</html>`
}

func formatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return "N/A"
	}
	// Convert Unix timestamp to time and format it
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

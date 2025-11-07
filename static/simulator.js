class BoidsSimulator {
    constructor() {
        this.canvas = document.getElementById('simulationCanvas');
        this.ctx = this.canvas.getContext('2d');
        this.previewObstacle = document.getElementById('previewObstacle');
        this.attractorElement = document.getElementById('attractor');
        
        // 获取画布的实际尺寸
        this.canvasWidth = this.canvas.width;
        this.canvasHeight = this.canvas.height;
        
        this.state = {
            birds: [],
            obstacles: [],
            attractor: { active: false, position: { x: 0, y: 0 } },
            step: 0,
            running: true
        };
        
        this.mode = 'none';
        this.isRunning = true;
        this.animationId = null;
        this.dragStart = null;
        this.dragEnd = null;
        this.isDragging = false;
        this.isAttracting = false;

        // 修复：初始化时立即获取状态
        this.initEventListeners();
        this.startSimulation();
        this.updateState(); // 立即获取初始状态
    }

    initEventListeners() {
        // 按钮事件
        document.getElementById('pauseBtn').addEventListener('click', () => this.togglePause());
        document.getElementById('resetBtn').addEventListener('click', () => this.reset());
        document.getElementById('addBirdBtn').addEventListener('click', () => this.setMode('addBird'));
        document.getElementById('addObstacleBtn').addEventListener('click', () => this.setMode('addObstacle'));
        document.getElementById('addAttractorBtn').addEventListener('click', () => this.setMode('attractor'));

        // 画布事件 - 使用正确的坐标转换
        this.canvas.addEventListener('click', (e) => this.handleCanvasClick(e));
        this.canvas.addEventListener('mousedown', (e) => this.handleMouseDown(e));
        this.canvas.addEventListener('mousemove', (e) => this.handleMouseMove(e));
        this.canvas.addEventListener('mouseup', (e) => this.handleMouseUp(e));
        this.canvas.addEventListener('mouseleave', (e) => this.handleMouseLeave(e));
        
        // 防止拖拽时选中文本
        this.canvas.addEventListener('selectstart', (e) => e.preventDefault());
    }

    // 修复：正确的坐标转换函数
    getCanvasCoordinates(e) {
        const rect = this.canvas.getBoundingClientRect();
        const scaleX = this.canvas.width / rect.width;    // 水平缩放比例
        const scaleY = this.canvas.height / rect.height;  // 垂直缩放比例
        
        return {
            x: (e.clientX - rect.left) * scaleX,
            y: (e.clientY - rect.top) * scaleY
        };
    }

    setMode(mode) {
        this.mode = mode;
        
        // 更新按钮状态
        document.getElementById('addBirdBtn').classList.toggle('active', mode === 'addBird');
        document.getElementById('addObstacleBtn').classList.toggle('active', mode === 'addObstacle');
        document.getElementById('addAttractorBtn').classList.toggle('active', mode === 'attractor');
        
        // 更新鼠标光标
        this.canvas.classList.toggle('cursor-crosshair', mode === 'addBird' || mode === 'addObstacle');
        this.canvas.classList.toggle('cursor-pointer', mode === 'attractor');
        
        let status = 'Ready';
        switch(mode) {
            case 'addBird': status = 'Click on canvas to add birds'; break;
            case 'addObstacle': status = 'Drag on canvas to create obstacles'; break;
            case 'attractor': status = 'Click and hold to attract birds'; break;
        }
        this.updateStatus(status);
    }

    handleCanvasClick(e) {
        if (this.mode === 'addBird') {
            const coords = this.getCanvasCoordinates(e);
            
            // 添加3-5只鸟，形成小群体
            const count = 3 + Math.floor(Math.random() * 3);
            for (let i = 0; i < count; i++) {
                setTimeout(() => {
                    const offsetX = (Math.random() - 0.5) * 40;
                    const offsetY = (Math.random() - 0.5) * 40;
                    this.addBird(coords.x + offsetX, coords.y + offsetY);
                }, i * 100);
            }
        }
    }

    handleMouseDown(e) {
        const coords = this.getCanvasCoordinates(e);
        
        if (this.mode === 'addObstacle') {
            this.dragStart = { x: coords.x, y: coords.y };
            this.isDragging = true;
            this.updatePreviewObstacle(coords.x, coords.y, 0);
        } else if (this.mode === 'attractor') {
            this.isAttracting = true;
            this.setAttractor(coords.x, coords.y, true);
        }
    }

    handleMouseMove(e) {
        const coords = this.getCanvasCoordinates(e);
        
        if (this.isDragging && this.dragStart) {
            const radius = Math.sqrt(
                Math.pow(coords.x - this.dragStart.x, 2) + 
                Math.pow(coords.y - this.dragStart.y, 2)
            );
            this.updatePreviewObstacle(this.dragStart.x, this.dragStart.y, radius);
        } else if (this.isAttracting) {
            this.setAttractor(coords.x, coords.y, true);
        }
    }

    handleMouseUp(e) {
        if (this.isDragging && this.dragStart) {
            const coords = this.getCanvasCoordinates(e);
            const radius = Math.sqrt(
                Math.pow(coords.x - this.dragStart.x, 2) + 
                Math.pow(coords.y - this.dragStart.y, 2)
            );
            
            if (radius > 5) { // 最小半径
                this.addObstacle(this.dragStart.x, this.dragStart.y, radius);
            }
            
            this.isDragging = false;
            this.dragStart = null;
            this.hidePreviewObstacle();
        } else if (this.isAttracting) {
            this.isAttracting = false;
            this.setAttractor(0, 0, false);
        }
    }

    handleMouseLeave(e) {
        if (this.isDragging) {
            this.isDragging = false;
            this.dragStart = null;
            this.hidePreviewObstacle();
        }
        if (this.isAttracting) {
            this.isAttracting = false;
            this.setAttractor(0, 0, false);
        }
    }

    updatePreviewObstacle(x, y, radius) {
        // 转换为CSS像素坐标用于显示
        const rect = this.canvas.getBoundingClientRect();
        const scaleX = rect.width / this.canvas.width;
        const scaleY = rect.height / this.canvas.height;
        
        this.previewObstacle.style.display = 'block';
        this.previewObstacle.style.left = (rect.left + (x - radius) * scaleX) + 'px';
        this.previewObstacle.style.top = (rect.top + (y - radius) * scaleY) + 'px';
        this.previewObstacle.style.width = (radius * 2 * scaleX) + 'px';
        this.previewObstacle.style.height = (radius * 2 * scaleY) + 'px';
    }

    hidePreviewObstacle() {
        this.previewObstacle.style.display = 'none';
    }

    updateAttractorDisplay(x, y, active) {
        if (active) {
            // 转换为CSS像素坐标用于显示
            const rect = this.canvas.getBoundingClientRect();
            const scaleX = rect.width / this.canvas.width;
            const scaleY = rect.height / this.canvas.height;
            
            this.attractorElement.style.display = 'block';
            this.attractorElement.style.left = (rect.left + (x - 10) * scaleX) + 'px';
            this.attractorElement.style.top = (rect.top + (y - 10) * scaleY) + 'px';
        } else {
            this.attractorElement.style.display = 'none';
        }
    }

    async addBird(x, y) {
        try {
            const response = await fetch('/api/bird', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ x, y })
            });
            await this.updateStateFromResponse(response);
            this.updateStatus(`Added bird at (${Math.round(x)}, ${Math.round(y)})`);
        } catch (error) {
            console.error('Error adding bird:', error);
        }
    }

    async addObstacle(x, y, radius) {
        try {
            const response = await fetch('/api/obstacle', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ x, y, radius })
            });
            await this.updateStateFromResponse(response);
            this.updateStatus(`Added obstacle (r=${Math.round(radius)})`);
        } catch (error) {
            console.error('Error adding obstacle:', error);
        }
    }

    async setAttractor(x, y, active) {
        try {
            const response = await fetch('/api/attractor', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ x, y, active })
            });
            await this.updateStateFromResponse(response);
            this.updateAttractorDisplay(x, y, active);
            if (active) {
                this.updateStatus('Attractor active - birds will follow your cursor');
            }
        } catch (error) {
            console.error('Error setting attractor:', error);
        }
    }

    async togglePause() {
        try {
            const response = await fetch('/api/toggle', { method: 'POST' });
            await this.updateStateFromResponse(response);
            
            const pauseBtn = document.getElementById('pauseBtn');
            if (this.state.running) {
                pauseBtn.textContent = '⏸️ Pause Simulation';
                this.updateStatus('Simulation running');
            } else {
                pauseBtn.textContent = '▶️ Continue Simulation';
                this.updateStatus('Simulation paused');
            }
        } catch (error) {
            console.error('Error toggling simulation:', error);
        }
    }

    async reset() {
        try {
            const response = await fetch('/api/reset', { method: 'POST' });
            if (response.ok) {
                // 修复：重置后立即获取新状态
                await this.updateState();
                this.updateStatus('Canvas reset');
                this.hidePreviewObstacle();
                this.updateAttractorDisplay(0, 0, false);
                // 修复：重置后重新开始模拟循环
                if (this.isRunning && !this.animationId) {
                    this.startSimulation();
                }
            }
        } catch (error) {
            console.error('Error resetting simulation:', error);
        }
    }

    // 新增：直接获取状态的方法
    async updateState() {
        try {
            const response = await fetch('/api/state');
            await this.updateStateFromResponse(response);
        } catch (error) {
            console.error('Error updating state:', error);
        }
    }

    async updateStateFromResponse(response) {
        if (response.ok) {
            const data = await response.json();
            this.state = data.state;
            this.updateUI();
        } else {
            throw new Error('Request failed');
        }
    }

    async stepSimulation() {
        try {
            const response = await fetch('/api/step', { method: 'POST' });
            await this.updateStateFromResponse(response);
        } catch (error) {
            console.error('Error stepping simulation:', error);
        }
    }

    startSimulation() {
        // 修复：确保只有一个动画循环在运行
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
        }
        
        const loop = async () => {
            if (this.state.running) {
                await this.stepSimulation();
                this.render();
            }
            this.animationId = requestAnimationFrame(loop.bind(this));
        };
        this.animationId = requestAnimationFrame(loop);
    }

    updateUI() {
        document.getElementById('birdCount').textContent = this.state.birds.length;
        document.getElementById('obstacleCount').textContent = this.state.obstacles.length;
        document.getElementById('stepCount').textContent = this.state.step;
        document.getElementById('statusText').textContent = this.state.running ? 'Running' : 'Paused';
    }

    updateStatus(message) {
        document.getElementById('status').textContent = message;
    }

    render() {
        // 修复：清空画布时清除所有内容，避免轨迹
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        // 绘制背景
        this.ctx.fillStyle = 'rgba(30, 60, 114, 1)'; // 改为不透明背景
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

        // 绘制障碍物
        this.ctx.fillStyle = 'rgba(200, 50, 50, 0.8)';
        this.ctx.strokeStyle = 'rgba(150, 30, 30, 0.9)';
        this.ctx.lineWidth = 2;
        
        this.state.obstacles.forEach(obstacle => {
            this.ctx.beginPath();
            this.ctx.arc(obstacle.position.x, obstacle.position.y, obstacle.radius, 0, Math.PI * 2);
            this.ctx.fill();
            this.ctx.stroke();
        });

        // 绘制鸟
        this.state.birds.forEach(bird => {
            this.drawBird(bird);
        });

        // 绘制引导点
        if (this.state.attractor.active) {
            this.ctx.fillStyle = 'rgba(255, 235, 59, 0.6)';
            this.ctx.beginPath();
            this.ctx.arc(this.state.attractor.position.x, this.state.attractor.position.y, 15, 0, Math.PI * 2);
            this.ctx.fill();
            
            this.ctx.strokeStyle = 'rgba(255, 193, 7, 0.8)';
            this.ctx.lineWidth = 2;
            this.ctx.stroke();
        }
    }

    drawBird(bird) {
        const x = bird.position.x;
        const y = bird.position.y;
        const dx = bird.velocity.dx;
        const dy = bird.velocity.dy;
        
        // 计算鸟的方向角度
        const angle = Math.atan2(dy, dx);
        
        this.ctx.save();
        this.ctx.translate(x, y);
        this.ctx.rotate(angle);
        
        // 鸟身体（三角形）
        this.ctx.fillStyle = `hsl(${(bird.id.charCodeAt(0) * 137) % 360}, 80%, 60%)`;
        this.ctx.strokeStyle = 'rgba(0, 0, 0, 0.5)';
        this.ctx.lineWidth = 1;
        
        this.ctx.beginPath();
        this.ctx.moveTo(6, 0);
        this.ctx.lineTo(-4, -3);
        this.ctx.lineTo(-4, 3);
        this.ctx.closePath();
        this.ctx.fill();
        this.ctx.stroke();
        
        // 鸟眼睛
        this.ctx.fillStyle = 'white';
        this.ctx.beginPath();
        this.ctx.arc(3, -1, 1, 0, Math.PI * 2);
        this.ctx.arc(3, 1, 1, 0, Math.PI * 2);
        this.ctx.fill();
        
        this.ctx.fillStyle = 'black';
        this.ctx.beginPath();
        this.ctx.arc(3.5, -1, 0.5, 0, Math.PI * 2);
        this.ctx.arc(3.5, 1, 0.5, 0, Math.PI * 2);
        this.ctx.fill();
        
        this.ctx.restore();
    }
}

// 页面加载完成后初始化模拟器
document.addEventListener('DOMContentLoaded', () => {
    new BoidsSimulator();
});
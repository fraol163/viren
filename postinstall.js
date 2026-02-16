const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

function log(msg) {
    console.log(`\x1b[96m${msg}\x1b[0m`);
}

function error(msg) {
    console.error(`\x1b[91mERROR: ${msg}\x1b[0m`);
    process.exit(1);
}

const isWindows = process.platform === 'win32';
const binDir = path.join(__dirname, 'bin');
const exePath = path.join(binDir, isWindows ? 'viren.exe' : 'viren');

if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
}

log('Building Viren binary...');

try {
    // Check if Go is installed
    execSync('go version', { stdio: 'ignore' });
} catch (e) {
    error('Go is required to build Viren. Please install Go from https://go.dev/dl/');
}

try {
    const buildCmd = `go build -ldflags "-s -w" -o "${exePath}" ./cmd/viren/main.go`;
    log(`Executing: ${buildCmd}`);
    execSync(buildCmd, { stdio: 'inherit', cwd: __dirname });
    log('Build successful!');
} catch (e) {
    error(`Failed to build Viren: ${e.message}`);
}

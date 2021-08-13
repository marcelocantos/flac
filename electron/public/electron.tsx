const path = require('path');

const { app, BrowserWindow, ipcMain } = require('electron');
const isDev = require('electron-is-dev');

const AsyncDB = require('./data/AsyncDB');
const Database = require('./data/Database');

const { default: installExtension, REACT_DEVELOPER_TOOLS } = require('electron-devtools-installer');

async function createWindow() {
  // Create the browser window.
  const win = new BrowserWindow({
    title: "flac",
    width: 800,
    height: 600,
    webPreferences: {
      contextIsolation: true,
      enableRemoteModule: false,
      nodeIntegration: false,
      webSecurity: true,
      preload: path.join(__dirname, "preload.js"), // use a preload script
    },
  });

  // and load the index.html of the app.
  // win.loadFile("index.html");
  win.loadURL(
    isDev
      ? 'http://localhost:3000'
      : `file://${path.join(__dirname, '../build/index.html')}`
  );
  // Open the DevTools.
  if (isDev) {
    win.webContents.openDevTools({ mode: 'detach' });
  }

  try {
    await installExtension(REACT_DEVELOPER_TOOLS);
    console.log(`Added Extension: ${name}`)
  } catch (error) {
    console.log('An error occurred:', error)
  }
}

ipcMain.handle("data", async (_, {words}) => {
  console.log(Database);
  const db = await AsyncDB.open('flac.db');
  const data = await Database.build(db, "", words);
  return {HeadWord: await data.HeadWord};
});

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(createWindow);

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  // if (process.platform !== 'darwin') {
    app.quit();
  // }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

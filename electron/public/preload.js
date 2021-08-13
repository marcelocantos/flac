const {
  contextBridge,
  ipcRenderer
} = require("electron");

// Expose protected methods that allow the renderer process to use
// the ipcRenderer without exposing the entire object
contextBridge.exposeInMainWorld(
  "api", {
      call: async (channel, ...args) => {
          // whitelist channels
          let validChannels = ["data"];
          if (validChannels.includes(channel)) {
              return await ipcRenderer.invoke(channel, ...args);
          } else {
              throw new Error(`Invalid channel name: "${channel}"`)
          }
      },
  }
);

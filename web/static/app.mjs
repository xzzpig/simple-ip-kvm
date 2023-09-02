import * as ws from './ws.mjs';
import * as kb from './kb.mjs';
import * as mouse from './mouse.mjs';

new Vue({
  el: '#app',
  data: {
    // serviceHostName: '10.0.0.235',
    serviceHostName: location.hostname,
    serviceHost: location.host,
    streamSrc: '',
    snapshotSrc: '',
    $channel: null,
    isKeyCaptureActive: false,
    isPointorLocked: false,
    mouseMoveSlice: [0, 0],
    activeDialog: '',
    pasteContent: '',
  },
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      var draggie = new Draggabilly( '.float-zone', {
        axis: 'x',
        containment: '#app'
      });
      try {
        const config = await this.fetchConfig();
        document.title = config.app_title;

        const streamOk = await this.pingStream(config.stream_port);
        if (!streamOk) {
          throw new Error(
            'Video stream is not ready, please check mjpeg process'
          );
        }
        this.$channel = await ws.init(
          location.protocol == 'http:' ? `ws://${this.serviceHost}/websocket` : `wss://${this.serviceHost}/websocket`
        );
        this.bindKeyHandler();
        this.bindMouseHandler();

        this.streamSrc = this.parseString(config.stream_url, config);
        this.snapshotSrc = this.parseString(config.snapshot_url, config);
      } catch (e) {
        alert(e.toString());
      }
    },
    parseString(str, config) {
      return str.replaceAll("{URL_PROTOCOL}", location.protocol).replaceAll("{URL_HOST}", location.hostname).replaceAll("{STREAM_PORT}", config.stream_port).replaceAll("{URL_HOSTPORT}",location.host)
    },
    async pingStream(port) {
      try {
        const pingRes = await fetch(this.snapshotSrc);
        return pingRes.status === 200;
      } catch (e) {
        return false;
      }
    },
    async fetchConfig() {
      try {
        const res = await fetch('/api/config');
        return res.json();
      } catch (e) {
        return null;
      }
    },
    bindKeyHandler() {
      document.addEventListener('keydown', (evt) => {
        if (!this.isKeyCaptureActive) {
          if (evt.key === 'Enter' && !this.activeDialog) {
            this.setScreenFocus(true);
          }
          return;
        }

        evt.preventDefault();

        if (evt.repeat) {
          return;
        }

        if (evt.key === 'Escape' && evt.shiftKey) {
          this.setScreenFocus(false);
          return;
        }
        kb.sendEvent(this.$channel, evt.key, 'keydown');
      });

      document.addEventListener('keyup', (evt) => {
        if (!this.isKeyCaptureActive) {
          return;
        }
        kb.sendEvent(this.$channel, evt.key, 'keyup');
      });
    },
    bindMouseHandler() {
      const mouseMoveSlice = this.mouseMoveSlice;

      document.addEventListener('pointerlockchange', (evt) => {
        this.isPointorLocked =
          document.pointerLockElement &&
          document.pointerLockElement.classList.contains('screen');
        mouse.sendEvent(this.$channel, '', 'reset');
      });

      window.setInterval(() => {
        if (mouseMoveSlice[0] !== 0 || mouseMoveSlice[1] !== 0) {
          mouse.sendEvent(this.$channel, mouseMoveSlice, 'move');
          mouseMoveSlice[0] = 0;
          mouseMoveSlice[1] = 0;
        }
      }, 30);

      mouse.sendEvent(this.$channel, 1, 'config-move-factor');
    },
    sendKey(key) {
      kb.sendEvent(this.$channel, key, 'keydown');
      setTimeout(() => {
        kb.sendEvent(this.$channel, key, 'keyup');
      }, 100)
    },
    onScreenBlur() {
      this.isKeyCaptureActive = false;
      if (this.isPointorLocked) {
        this.setPointerLock(false);
      }
      kb.sendEvent(this.$channel, '', 'reset');
    },
    onScreenFocus() {
      this.setDialog();
      this.isKeyCaptureActive = true;
      kb.sendEvent(this.$channel, '', 'reset');
    },
    setScreenFocus(bool) {
      const screen = document.querySelector('.screen');
      screen[bool ? 'focus' : 'blur']();
    },
    setPointerLock(bool) {
      const screen = document.querySelector('.screen');
      if (bool) {
        try {
          this.setDialog();
          screen.requestPointerLock();
        } catch (e) { }
      } else {
        document.exitPointerLock();
      }
    },
    onScreenMouseMove(evt) {
      if (!this.isPointorLocked) {
        return;
      }
      this.mouseMoveSlice[0] += evt.movementX;
      this.mouseMoveSlice[1] += evt.movementY;
    },
    onScreenMouseDown(evt) {
      if (!this.isPointorLocked) {
        if (evt.button === 0) {
          this.setPointerLock(true);
        }
        return;
      }
      evt.preventDefault();
      mouse.sendEvent(this.$channel, evt.button, 'mousedown');
    },
    onScreenMouseUp(evt) {
      if (!this.isPointorLocked) {
        return;
      }
      mouse.sendEvent(this.$channel, evt.button, 'mouseup');
    },
    onScreenMouseWheel(evt) {
      if (!this.isPointorLocked) {
        return;
      }
      mouse.sendEvent(this.$channel, evt.wheelDeltaY, 'wheel');
    },
    doRemotePaste() {
      kb.sendSequence(this.$channel, this.pasteContent);
      this.pasteContent = '';
    },
    setDialog(name) {
      if (name) {
        this.setPointerLock(false);
        this.setScreenFocus(false);
        this.activeDialog = name;
      } else {
        this.activeDialog = '';
      }
    },
  },
});

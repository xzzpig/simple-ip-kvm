#include "HID-Project.h"

#define KB_EVT_START 248
#define MOUSE_EVT_START 249
#define KEY_SEQUENCE_EVT_START 250
#define EVT_END 251

#define KB_EVT_TYPE_KEYDOWN 1
#define KB_EVT_TYPE_KEYUP 2
#define KB_EVT_TYPE_RESET 3

#define MOUSE_EVT_TYPE_MOVE 1
#define MOUSE_EVT_TYPE_LEFT_DOWN 2
#define MOUSE_EVT_TYPE_LEFT_UP 3
#define MOUSE_EVT_TYPE_MIDDLE_DOWN 4
#define MOUSE_EVT_TYPE_MIDDLE_UP 5
#define MOUSE_EVT_TYPE_RIGHT_DOWN 6
#define MOUSE_EVT_TYPE_RIGHT_UP 7
#define MOUSE_EVT_TYPE_WHEEL 8
#define MOUSE_EVT_TYPE_RESET 9
#define MOUSE_EVT_TYPE_CONFIG_MOVE_FACTOR 10

#define SERIAL_KEY_CONTROL 0x80
#define SERIAL_KEY_SHIFT 0x81
#define SERIAL_KEY_ALT 0x82
#define SERIAL_KEY_META 0x83
#define SERIAL_KEY_TAB 0xB3
#define SERIAL_KEY_CAPSLOCK 0xC1
#define SERIAL_KEY_BACKSPACE 0xB2
#define SERIAL_KEY_ENTER 0xB0
#define SERIAL_KEY_CONTEXTMENU 0xED
#define SERIAL_KEY_INSERT 0xD1
#define SERIAL_KEY_DELETE 0xD4
#define SERIAL_KEY_HOME 0xD2
#define SERIAL_KEY_END 0xD5
#define SERIAL_KEY_PAGEUP 0xD3
#define SERIAL_KEY_PAGEDOWN 0xD6
#define SERIAL_KEY_ARROWUP 0xDA
#define SERIAL_KEY_ARROWDOWN 0xD9
#define SERIAL_KEY_ARROWLEFT 0xD8
#define SERIAL_KEY_ARROWRIGHT 0xD7
#define SERIAL_KEY_PRINTSCREEN 0xCE
#define SERIAL_KEY_SCROLLLOCK 0xCF
#define SERIAL_KEY_PAUSE 0xD0
#define SERIAL_KEY_ESCAPE 0xB1
#define SERIAL_KEY_F01 0xC2
#define SERIAL_KEY_F12 0xCD

#define R_BUF_LEN 32

bool led = false;

int rBuf[R_BUF_LEN];
int rBufCursor = 0;
int mouseMoveFactor = 1;

void blink() {
  digitalWrite(LED_BUILTIN, led ? HIGH : LOW);
  led = !led;
}

void setup() {
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, LOW);

  BootKeyboard.begin();
  Mouse.begin();

  Serial1.begin(19200);
}

void kbCodeAction(int type, KeyboardKeycode keyCode) {
  if (type == KB_EVT_TYPE_KEYDOWN) BootKeyboard.press(keyCode);
  else BootKeyboard.release(keyCode);
}

void kbAction(int type, int keyCode) {
  if (keyCode >= 32 && keyCode <= 126) {  //asii
    if (type == KB_EVT_TYPE_KEYDOWN) BootKeyboard.press(keyCode);
    else BootKeyboard.release(keyCode);
  }
  if (keyCode >= SERIAL_KEY_F01 && keyCode <= SERIAL_KEY_F12) {  //F1-F12
    kbCodeAction(type, (KeyboardKeycode)(keyCode - SERIAL_KEY_F01 + KEY_F1));
  }
  switch (keyCode) {
    case SERIAL_KEY_CONTROL:
      kbCodeAction(type, KEY_LEFT_CTRL);
      break;
    case SERIAL_KEY_SHIFT:
      kbCodeAction(type, KEY_LEFT_SHIFT);
      break;
    case SERIAL_KEY_ALT:
      kbCodeAction(type, KEY_LEFT_ALT);
      break;
    case SERIAL_KEY_META:
      kbCodeAction(type, KEY_LEFT_GUI);
      break;
    case SERIAL_KEY_TAB:
      kbCodeAction(type, KEY_TAB);
      break;
    case SERIAL_KEY_CAPSLOCK:
      kbCodeAction(type, KEY_CAPS_LOCK);
      break;
    case SERIAL_KEY_BACKSPACE:
      kbCodeAction(type, KEY_BACKSPACE);
      break;
    case SERIAL_KEY_ENTER:
      kbCodeAction(type, KEY_ENTER);
      break;
    case SERIAL_KEY_CONTEXTMENU:
      kbCodeAction(type, KEY_MENU);
      break;
    case SERIAL_KEY_INSERT:
      kbCodeAction(type, KEY_INSERT);
      break;
    case SERIAL_KEY_DELETE:
      kbCodeAction(type, KEY_DELETE);
      break;
    case SERIAL_KEY_HOME:
      kbCodeAction(type, KEY_HOME);
      break;
    case SERIAL_KEY_END:
      kbCodeAction(type, KEY_END);
      break;
    case SERIAL_KEY_PAGEUP:
      kbCodeAction(type, KEY_PAGE_UP);
      break;
    case SERIAL_KEY_PAGEDOWN:
      kbCodeAction(type, KEY_PAGE_DOWN);
      break;
    case SERIAL_KEY_ARROWUP:
      kbCodeAction(type, KEY_UP_ARROW);
      break;
    case SERIAL_KEY_ARROWDOWN:
      kbCodeAction(type, KEY_DOWN_ARROW);
      break;
    case SERIAL_KEY_ARROWLEFT:
      kbCodeAction(type, KEY_LEFT_ARROW);
      break;
    case SERIAL_KEY_ARROWRIGHT:
      kbCodeAction(type, KEY_RIGHT_ARROW);
      break;
    case SERIAL_KEY_PRINTSCREEN:
      kbCodeAction(type, KEY_PRINTSCREEN);
      break;
    case SERIAL_KEY_SCROLLLOCK:
      kbCodeAction(type, KEY_SCROLL_LOCK);
      break;
    case SERIAL_KEY_PAUSE:
      kbCodeAction(type, KEY_PAUSE);
      break;
    case SERIAL_KEY_ESCAPE:
      kbCodeAction(type, KEY_ESC);
      break;
    default:
      if (type == KB_EVT_TYPE_KEYDOWN) BootKeyboard.press(keyCode);
      else BootKeyboard.release(keyCode);
      break;
  }
}

void parse_r_buf() {
  if (rBuf[0] == KB_EVT_START && rBufCursor == 3) {
    switch (rBuf[1]) {
      case KB_EVT_TYPE_KEYDOWN:
        kbAction(KB_EVT_TYPE_KEYDOWN, rBuf[2]);
        break;
      case KB_EVT_TYPE_KEYUP:
        kbAction(KB_EVT_TYPE_KEYUP, rBuf[2]);
        break;
      case KB_EVT_TYPE_RESET:
        BootKeyboard.releaseAll();
        break;
    }
  }

  if (rBuf[0] == MOUSE_EVT_START && rBufCursor == 4) {
    switch (rBuf[1]) {
      case MOUSE_EVT_TYPE_MOVE:
        Mouse.move(mouseMoveFactor * (rBuf[2] - 120), mouseMoveFactor * (rBuf[3] - 120), 0);
        break;
      case MOUSE_EVT_TYPE_LEFT_DOWN:
        Mouse.press(MOUSE_LEFT);
        break;
      case MOUSE_EVT_TYPE_LEFT_UP:
        Mouse.release(MOUSE_LEFT);
        break;
      case MOUSE_EVT_TYPE_RIGHT_DOWN:
        Mouse.press(MOUSE_RIGHT);
        break;
      case MOUSE_EVT_TYPE_RIGHT_UP:
        Mouse.release(MOUSE_RIGHT);
        break;
      case MOUSE_EVT_TYPE_MIDDLE_DOWN:
        Mouse.press(MOUSE_MIDDLE);
        break;
      case MOUSE_EVT_TYPE_MIDDLE_UP:
        Mouse.release(MOUSE_MIDDLE);
        break;
      case MOUSE_EVT_TYPE_WHEEL:
        Mouse.move(0, 0, rBuf[2] - 120);
        break;
      case MOUSE_EVT_TYPE_RESET:
        Mouse.release(MOUSE_LEFT);
        Mouse.release(MOUSE_RIGHT);
        Mouse.release(MOUSE_MIDDLE);
        break;
      case MOUSE_EVT_TYPE_CONFIG_MOVE_FACTOR:
        mouseMoveFactor = rBuf[2];
    }
  }

  if (rBuf[0] == KEY_SEQUENCE_EVT_START && rBufCursor > 1) {
    BootKeyboard.releaseAll();
    for (int i = 1; i < rBufCursor; i += 1) {
      BootKeyboard.write(rBuf[i]);
    }
  }
}

void reset_r_buf() {
  rBufCursor = 0;
  rBuf[0] = 0;
}

void loop() {
  int curVal;
  while (Serial1.available() > 0) {
    curVal = Serial1.read();

    if (curVal == EVT_END) {
      parse_r_buf();
      blink();
      reset_r_buf();
    } else {
      if (rBufCursor == 0) {
        if (curVal == KB_EVT_START || curVal == MOUSE_EVT_START || curVal == KEY_SEQUENCE_EVT_START) {
          rBuf[rBufCursor] = curVal;
          rBufCursor += 1;
        }
      } else {
        if (rBufCursor < R_BUF_LEN) {
          rBuf[rBufCursor] = curVal;
          rBufCursor += 1;
        } else {
          // overflow, reset rBuf
          rBuf[0] = 0;
        }
      }
    }
  }
}

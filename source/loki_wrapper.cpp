#include "loki_wrapper.h"
#include <stdexcept>
#include <stdio.h>
#include <stdlib.h>

#include "mouse.h"
#include "keyboard.h"

// int count = 0;
// WindowInfo windows[1024];

// BOOL CALLBACK enumChildWindowsProc(HWND hwnd, LPARAM lParam)
// {
//   WindowInfo *parent = (WindowInfo *)lParam;
//   WindowInfo *childWi = (WindowInfo *)malloc(sizeof(WindowInfo));

//   childWi->hwnd = hwnd;
//   GetWindowRect(hwnd, &childWi->rect);
//   GetClientRect(hwnd, &childWi->clientRect);
//   childWi->style = GetWindowLong(hwnd, GWL_STYLE);
//   childWi->exStyle = GetWindowLong(hwnd, GWL_EXSTYLE);
//   GetWindowText(hwnd, childWi->title, sizeof(childWi->title));
//   GetClassName(hwnd, childWi->class_name, sizeof(childWi->class_name));
//   childWi->childCount = 0; // 子窗口不再进一步枚举

//   if (parent->childCount < MAX_CHILDREN)
//   {
//     parent->children[parent->childCount++] = childWi;
//   }

//   return TRUE;
// }

// void copyWindowInfo(WindowInfo *dest, const WindowInfo *src)
// {
//   *dest = *src;
//   for (int i = 0; i < src->childCount; ++i)
//   {
//     dest->children[i] = (WindowInfo *)malloc(sizeof(WindowInfo));
//     *dest->children[i] = *src->children[i];
//   }
// }

// BOOL CALLBACK enumWindowsProc(HWND hwnd, LPARAM lParam)
// {
//   DWORD windowPID;
//   GetWindowThreadProcessId(hwnd, &windowPID);
//   if (windowPID != *(DWORD *)lParam)
//   {
//     return TRUE;
//   }

//   WindowInfo wi;
//   wi.hwnd = hwnd;
//   GetWindowRect(hwnd, &wi.rect);
//   GetClientRect(hwnd, &wi.clientRect);
//   wi.style = GetWindowLong(hwnd, GWL_STYLE);
//   wi.exStyle = GetWindowLong(hwnd, GWL_EXSTYLE);
//   GetWindowText(hwnd, wi.title, sizeof(wi.title));
//   GetClassName(hwnd, wi.class_name, sizeof(wi.class_name));
//   wi.childCount = 0;

//   EnumChildWindows(hwnd, enumChildWindowsProc, (LPARAM)&wi);

//   if (count < 1024)
//   {
//     copyWindowInfo(&windows[count], &wi);
//     count++;
//   }

//   for (int i = 0; i < wi.childCount; ++i)
//   {
//     free(wi.children[i]);
//   }

//   return TRUE;
// }

// void enumerateWindows(DWORD pid)
// {
//   count = 0;
//   EnumWindows(enumWindowsProc, (LPARAM)&pid);
// }

void *new_mouse()
{
  return new Mouse();
}

int mouse_initialize(void *mouse)
{
  try
  {
    static_cast<Mouse *>(mouse)->initialize();
    return 0; // Success
  }
  catch (const std::runtime_error &e)
  {
    // Log error or handle it as needed
    return 1; // Failure
  }
}

void mouse_move_cursor(void *mouse, int x, int y)
{
  static_cast<Mouse *>(mouse)->moveCursor(x, y);
}

void mouse_left_button_click(void *mouse)
{
  static_cast<Mouse *>(mouse)->leftButtonClick();
}

void delete_mouse(void *mouse)
{
  delete static_cast<Mouse *>(mouse);
}

void *new_keyboard()
{
  return new Keyboard();
}

int keyboard_initialize(void *keyboard)
{
  try
  {
    static_cast<Keyboard *>(keyboard)->initialize();
    return 0; // Success
  }
  catch (const std::runtime_error &e)
  {
    // Handle or log the exception as needed
    return 1; // Failure
  }
}

void keyboard_type(void *keyboard, unsigned char keyCode)
{
  Keyboard *kb = static_cast<Keyboard *>(keyboard);
  kb->type(keyCode);
}

void keyboard_key_down(void *keyboard, unsigned char keyCode)
{
  static_cast<Keyboard *>(keyboard)->keyDown(keyCode);
}

void keyboard_key_up(void *keyboard, unsigned char keyCode)
{
  static_cast<Keyboard *>(keyboard)->keyUp(keyCode);
}

void delete_keyboard(void *keyboard)
{
  delete static_cast<Keyboard *>(keyboard);
}

#ifdef __cplusplus
extern "C" {
#endif

void* new_mouse();
int mouse_initialize(void* mouse);
void mouse_move_cursor(void* mouse, int x, int y);
void mouse_left_button_click(void* mouse);
void delete_mouse(void* mouse);

void* new_keyboard();
int keyboard_initialize(void* keyboard);
void keyboard_type(void *keyboard, unsigned char keyCode);
void keyboard_key_down(void *keyboard, unsigned char keyCode);
void keyboard_key_up(void *keyboard, unsigned char keyCode);
void delete_keyboard(void* keyboard);

#ifdef __cplusplus
}
#endif

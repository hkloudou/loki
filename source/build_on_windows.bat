g++ -m32 -c -std=c++11 -Iinclude -ooutput/device.o source/device.cpp
g++ -m32 -c -std=c++11 -Iinclude -ooutput/mouse.o source/mouse.cpp
g++ -m32 -c -std=c++11 -Iinclude -ooutput/keyboard.o source/keyboard.cpp
g++ -m32 -c -std=c++11 -Iinclude -ooutput/registry.o source/registry.cpp
g++ -m32 -c -std=c++11 -Iinclude -ooutput/loki_wrapper.o source/loki_wrapper.cpp

ar rcs lib/32/libloki.a output/mouse.o output/keyboard.o output/device.o output/registry.o output/loki_wrapper.o


g++ -m64 -c -std=c++11 -Iinclude -ooutput64/device.o source/device.cpp
g++ -m64 -c -std=c++11 -Iinclude -ooutput64/mouse.o source/mouse.cpp
g++ -m64 -c -std=c++11 -Iinclude -ooutput64/keyboard.o source/keyboard.cpp
g++ -m64 -c -std=c++11 -Iinclude -ooutput64/registry.o source/registry.cpp
g++ -m64 -c -std=c++11 -Iinclude -ooutput64/loki_wrapper.o source/loki_wrapper.cpp

ar rcs lib/64/libloki.a output64/mouse.o output64/keyboard.o output64/device.o output64/registry.o output64/loki_wrapper.o

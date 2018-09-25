all:
	$(MAKE) -C demo/2/_mymath
	$(MAKE) -C demo/4/_luajit

clean:
	$(MAKE) -C demo/2/_mymath clean
	$(MAKE) -C demo/4/_luajit clean

.PHONY: all clean

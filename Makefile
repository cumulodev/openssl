NAME = openssl-1.0.2a
SRC = https://www.openssl.org/source/$(NAME).tar.gz
SRCDIR = vendor/$(NAME)
PREFIX ?= $(shell pwd)/build
OBJ = $(PREFIX)/lib/libcrypto.a $(PREFIX)/lib/libssl.a

export PKG_CONFIG_PATH = prebuilt/lib/pkgconfig

.PHONY: all install clean

all: $(SRCDIR) $(OBJ)

$(OBJ):
	mkdir -p $(PREFIX)
	cd $(SRCDIR) && ./config --prefix=$(PREFIX) --openssldir=$(PREFIX) no-shared
	make -C $(SRCDIR) depend
	make -C $(SRCDIR)
	make -C $(SRCDIR) install 

install: $(OBJ)
	go install -a

$(SRCDIR):
	wget $(SRC) -O vendor/$(NAME).tar.gz
	wget $(SRC).asc -O vendor/$(NAME).tar.gz.asc
	gpg vendor/$(NAME).tar.gz.asc
	tar -C vendor -x -z -f vendor/$(NAME).tar.gz

clean:
	rm -f vendor/$(NAME).tar.gz
	rm -f vendor/$(NAME).tar.gz.asc
	rm -rf $(SRCDIR)
	rm -rf $(PREFIX)

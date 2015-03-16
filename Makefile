CWD = $(shell pwd)
SSLDIR = ${CWD}/vendor/openssl-0.9.8ze
BUILDDIR = ${CWD}/out
OBJ = ${BUILDDIR}/lib/libcrypto.a ${BUILDDIR}/lib/libssl.a

export PKG_CONFIG_PATH = ${BUILDDIR}/lib/pkgconfig

.PHONY: all install clean

all: ${OBJ}

${OBJ}:
	mkdir -p ${BUILDDIR}
	cd ${SSLDIR} && ./config --prefix=${BUILDDIR} --openssldir=${BUILDDIR}
	make -C ${SSLDIR}
	make -C ${SSLDIR} install 

install: ${OBJ}
	go install

clean:
	make -C ${SSLDIR} clean
	rm -r ${BUILDDIR}
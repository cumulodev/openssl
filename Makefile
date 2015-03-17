CWD = $(shell pwd)
SSLDIR = ${CWD}/vendor/openssl-1.0.2
BUILDDIR = ${CWD}/out
OBJ = ${BUILDDIR}/lib/libcrypto.a ${BUILDDIR}/lib/libssl.a

export PKG_CONFIG_PATH = ${BUILDDIR}/lib/pkgconfig

.PHONY: all install clean

all: ${OBJ}

${OBJ}:
	mkdir -p ${BUILDDIR}
	cd ${SSLDIR} && ./config --prefix=${BUILDDIR} --openssldir=${BUILDDIR} --no-shared
	make -C ${SSLDIR} depend
	make -C ${SSLDIR}
	make -C ${SSLDIR} install 

install: ${OBJ}
	go install

clean:
	make -C ${SSLDIR} clean
	rm -r ${BUILDDIR}

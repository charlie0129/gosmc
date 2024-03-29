/*
 * Apple System Management Control (SMC) Tool
 * Copyright (C) 2006 devnull
 * Portions Copyright (C) 2013 Michael Wilber
 * Portions Copyright (C) 2023 Charlie Chiang
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.

 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

#ifndef __SMC_H__
#define __SMC_H__

#include <IOKit/IOKitLib.h>

#if (MAC_OS_X_VERSION_MIN_REQUIRED < 120000) // Before macOS 12 Monterey
       #define  IOMasterPort
#endif

#define OP_NONE 0
#define OP_LIST 1
#define OP_READ 2
#define OP_READ_FAN 3
#define OP_WRITE 4
#define OP_READ_TEMPS 5

#define KERNEL_INDEX_SMC 2

#define SMC_CMD_READ_BYTES 5
#define SMC_CMD_WRITE_BYTES 6
#define SMC_CMD_READ_INDEX 8
#define SMC_CMD_READ_KEYINFO 9
#define SMC_CMD_READ_PLIMIT 11
#define SMC_CMD_READ_VERS 12

typedef struct {
    char major;
    char minor;
    char build;
    char reserved[1];
    UInt16 release;
} SMCKeyData_vers_t;

typedef struct {
    UInt16 version;
    UInt16 length;
    UInt32 cpuPLimit;
    UInt32 gpuPLimit;
    UInt32 memPLimit;
} SMCKeyData_pLimitData_t;

typedef struct {
    UInt32 dataSize;
    UInt32 dataType;
    char dataAttributes;
} SMCKeyData_keyInfo_t;

typedef unsigned char SMCBytes_t[32];

typedef struct {
    UInt32 key;
    SMCKeyData_vers_t vers;
    SMCKeyData_pLimitData_t pLimitData;
    SMCKeyData_keyInfo_t keyInfo;
    char result;
    char status;
    char data8;
    UInt32 data32;
    SMCBytes_t bytes;
} SMCKeyData_t;

typedef char UInt32Char_t[5];

typedef struct {
    UInt32Char_t key;
    UInt32 dataSize;
    UInt32Char_t dataType;
    SMCBytes_t bytes;
} SMCVal_t;

kern_return_t SMCReadKey(UInt32Char_t key, SMCVal_t* val);
kern_return_t SMCWriteSimple(char* key, unsigned char* bytes, int len, io_connect_t conn);

kern_return_t SMCOpen(io_connect_t* conn);
kern_return_t SMCClose(io_connect_t conn);
kern_return_t SMCReadKey2(UInt32Char_t key, SMCVal_t* val, io_connect_t conn);

#endif

INFO  Using log config file: jar:file:/opt/ghidra/Ghidra/Framework/Generic/lib/Generic.jar!/generic.log4j.xml (LoggingInitialization)  
INFO  Using log file: /work/home/.config/ghidra/ghidra_12.1.4_PUBLIC/application.log (LoggingInitialization)  
INFO  Loading user preferences: /work/home/.config/ghidra/ghidra_12.1.4_PUBLIC/preferences (Preferences)  
INFO  Searching for classes... (ClassSearcher)  
INFO  Class search complete (620 ms) (ClassSearcher)  
INFO  Initializing SSL Context (DefaultSSLContextInitializer)  
INFO  Initializing Random Number Generator... (SecureRandomFactory)  
INFO  Random Number Generator initialization complete: NativePRNGNonBlocking (SecureRandomFactory)  
INFO  Trust manager disabled, cacerts have not been set (DefaultTrustManagerFactory)  
INFO  Starting cache cleanup: /var/tmp/ubuntu-ghidra/fscache2 (FileCacheMaintenanceDaemon)  
INFO  Finished cache cleanup, estimated storage used: 0 (FileCacheMaintenanceDaemon)  
INFO  Headless startup complete (1771 ms) (AnalyzeHeadless)  
INFO  Class searcher loaded 59 extension points (20 false positives) (ClassSearcher)  
INFO  HEADLESS Script Paths:
    /opt/ghidra/Ghidra/Features/SystemEmulation/ghidra_scripts
    /opt/ghidra/Ghidra/Features/GnuDemangler/ghidra_scripts
    /opt/ghidra/Ghidra/Features/WildcardAssembler/ghidra_scripts
    /opt/ghidra/Ghidra/Processors/DATA/ghidra_scripts
    /opt/ghidra/Ghidra/Features/DecompilerDependent/ghidra_scripts
    /opt/ghidra/Ghidra/Processors/JVM/ghidra_scripts
    /opt/ghidra/Ghidra/Debug/Debugger-rmi-trace/ghidra_scripts
    /opt/ghidra/Ghidra/Features/MicrosoftCodeAnalyzer/ghidra_scripts
    /opt/ghidra/Ghidra/Features/VersionTracking/ghidra_scripts
    /opt/ghidra/Ghidra/Features/PyGhidra/ghidra_scripts
    /opt/ghidra/Ghidra/Features/SwiftDemangler/ghidra_scripts
    /opt/ghidra/Ghidra/Features/Base/ghidra_scripts
    /opt/ghidra/Ghidra/Processors/Atmel/ghidra_scripts
    /opt/ghidra/Ghidra/Debug/Debugger/ghidra_scripts
    /opt/ghidra/Ghidra/Features/Decompiler/ghidra_scripts
    /opt/ghidra/Ghidra/Features/FunctionID/ghidra_scripts
    /opt/ghidra/Ghidra/Features/PDB/ghidra_scripts
    /opt/ghidra/Ghidra/Processors/8051/ghidra_scripts
    /opt/ghidra/Ghidra/Processors/PIC/ghidra_scripts
    /opt/ghidra/Ghidra/Features/FileFormats/ghidra_scripts
    /opt/ghidra/Ghidra/Features/BSim/ghidra_scripts
    /opt/ghidra/Ghidra/Features/BytePatterns/ghidra_scripts
    /work/scripts (HeadlessAnalyzer)  
INFO  HEADLESS: execution starts (HeadlessAnalyzer)  
INFO  Opening existing project: /work/projects/VoiceText-port-stage27 (HeadlessAnalyzer)  
INFO  Opening project: /work/projects/VoiceText-port-stage27 (HeadlessProject)  
INFO  REPORT: Processing project file: /vt_pau.dll (HeadlessAnalyzer)  
INFO  REPORT: Execute script: DecompileNamedFunctions.java '0x10017030' '0x10017090' '0x10017100' '0x100560a0'  (HeadlessAnalyzer)  
INFO  SCRIPT: /work/scripts/DecompileNamedFunctions.java (HeadlessAnalyzer)  
INFO  DecompileNamedFunctions.java> ===== 0x10017030 ===== (GhidraScript)  
INFO  DecompileNamedFunctions.java> Function: FUN_10017030 @ 10017030 (GhidraScript)  
INFO  DecompileNamedFunctions.java> 
uint __cdecl FUN_10017030(int param_1,int param_2)

{
  char cVar1;
  
  if (-1 < param_2) {
    cVar1 = *(char *)(param_1 + 0x2e8 + param_2);
    if ((&DAT_1007b9e0)[cVar1] == '\x01') {
      return CONCAT22((short)((uint)param_2 >> 0x10),(ushort)(byte)(&DAT_1007baa8)[cVar1]);
    }
    if ((cVar1 == '6') && (0 < param_2)) {
      cVar1 = *(char *)(param_1 + 0x2e7 + param_2);
      param_2 = (int)cVar1;
      if ((&DAT_1007b9e0)[param_2] == '\x01') {
        return CONCAT22(cVar1 >> 7,(ushort)(byte)(&DAT_1007baa8)[param_2]);
      }
    }
  }
  return CONCAT22((short)((uint)param_2 >> 0x10),0xffff);
}

 (GhidraScript)  
INFO  DecompileNamedFunctions.java> ===== 0x10017090 ===== (GhidraScript)  
INFO  DecompileNamedFunctions.java> Function: FUN_10017090 @ 10017090 (GhidraScript)  
INFO  DecompileNamedFunctions.java> 
short __cdecl FUN_10017090(int param_1,int param_2)

{
  char cVar1;
  
  if (param_2 < (int)(uint)*(byte *)(param_1 + 0x95)) {
    cVar1 = *(char *)(param_1 + 0x2e8 + param_2);
    if ((&DAT_1007b9e0)[cVar1] == '\x01') {
      return (short)(char)(&DAT_1007bb70)[cVar1];
    }
    if (((cVar1 == 'C') && (param_2 < (int)(*(byte *)(param_1 + 0x95) - 1))) &&
       (cVar1 = *(char *)(param_1 + 0x2e9 + param_2), (&DAT_1007b9e0)[cVar1] == '\x01')) {
      return (short)(char)(&DAT_1007bb70)[cVar1];
    }
  }
  return -1;
}

 (GhidraScript)  
INFO  DecompileNamedFunctions.java> ===== 0x10017100 ===== (GhidraScript)  
INFO  DecompileNamedFunctions.java> Function: FUN_10017100 @ 10017100 (GhidraScript)  
INFO  DecompileNamedFunctions.java> 
undefined4 __cdecl FUN_10017100(int param_1)

{
  byte bVar1;
  char cVar2;
  int iVar3;
  short sVar4;
  int iVar5;
  uint uVar6;
  int iVar7;
  byte *pbVar8;
  byte *pbVar9;
  int iVar10;
  byte *pbVar11;
  int iVar12;
  int local_1c;
  short local_14;
  int local_c;
  
  iVar3 = *(int *)(param_1 + 0x4c);
  local_1c = 0;
  if (0 < *(short *)(iVar3 + 2)) {
    pbVar11 = (byte *)(iVar3 + 0x6e1);
    local_c = -0x934 - iVar3;
    do {
      iVar5 = 0;
      if (0 < (int)(*pbVar11 - 1)) {
        do {
          if (((pbVar11[iVar5 + 0x294] == 0x30) &&
              ((&DAT_1007be98)[(char)pbVar11[iVar5 + 0x253]] != 0)) &&
             (pbVar11[iVar5 + 0x254] == 0x36)) {
            pbVar11[iVar5 + 0x253] = (&DAT_1007be98)[(char)pbVar11[iVar5 + 0x253]];
          }
          iVar5 = iVar5 + 1;
        } while (iVar5 < (int)(*pbVar11 - 1));
      }
      param_1 = 1;
      if (1 < (int)(*pbVar11 - 1)) {
        pbVar9 = pbVar11 + 0x254;
        pbVar8 = pbVar9;
        do {
          if (pbVar11[param_1 + 0x293] == 0x30) {
            uVar6 = FUN_10017030((int)(pbVar11 + -0x95),(int)(pbVar8 + local_c + 1 + -2));
            sVar4 = FUN_10017090((int)(pbVar11 + -0x95),(int)(pbVar8 + local_c + 1));
            local_14 = (short)uVar6;
            if (pbVar11[param_1 + 0x294] == 0x30) {
              bVar1 = *pbVar8;
              if (bVar1 == 0x39) {
                bVar1 = pbVar11[param_1 + 0x252];
                if (((bVar1 == 0x37) || (bVar1 == 0x38)) && (-1 < sVar4)) {
                  *pbVar8 = 0x48;
                }
                else if (((&DAT_1007be3c)[(char)bVar1] == '\x01') && (sVar4 == 0)) {
                  *pbVar8 = 0x48;
                }
                else if ((local_14 < 1) ||
                        (((((int)(*pbVar11 - 2) <= param_1 || (pbVar11[param_1 + 0x295] != 0x30)) ||
                          (pbVar9[param_1] != 7)) || (pbVar11[param_1 + 0x255] != 0x2d)))) {
                  if (-1 < local_14) {
                    if (sVar4 == 0) {
                      *pbVar8 = 0x4d;
                    }
                    else {
                      bVar1 = pbVar9[param_1];
                      if (((bVar1 == 0x2b) || (bVar1 == 0x2c)) || (bVar1 == 0x2d)) {
                        *pbVar8 = 0x4b;
                      }
                    }
                  }
                }
                else {
                  *pbVar8 = 0x4b;
                }
              }
              else if ((bVar1 == 0x2a) || (bVar1 == 0x35)) {
                if ((-1 < sVar4) && ((pbVar11[param_1 + 0x252] == 0x37 || (sVar4 == 0)))) {
                  if (bVar1 == 0x2a) {
                    *pbVar8 = 0x46;
                  }
                  else if (bVar1 == 0x35) {
                    *pbVar8 = 0x47;
                  }
                }
              }
              else if (((bVar1 == 0x15) && (-1 < local_14)) && (sVar4 == 0)) {
                *pbVar8 = 0x4e;
              }
            }
            else {
              if ((-1 < local_14) && (-1 < sVar4)) {
                if (pbVar9[param_1] != 0x43) {
                  if (*pbVar8 == 0x39) {
                    *pbVar8 = 0x4d;
                  }
                  else if (*pbVar8 == 0x15) {
                    *pbVar8 = 0x4e;
                  }
                }
                if (*pbVar8 == 0x2a) {
                  *pbVar8 = 0x46;
                }
                else if (*pbVar8 == 0x35) {
                  *pbVar8 = 0x47;
                }
              }
              if (((*pbVar8 == 0x39) && ((&DAT_1007be3c)[(char)pbVar11[param_1 + 0x252]] == '\x01'))
                 && (-1 < sVar4)) {
                *pbVar8 = 0x48;
              }
            }
          }
          param_1 = param_1 + 1;
          pbVar8 = pbVar8 + 1;
        } while (param_1 < (int)(*pbVar11 - 1));
      }
      iVar5 = 0;
      if (*pbVar11 != 1 && -1 < (int)(*pbVar11 - 1)) {
        pbVar9 = pbVar11 + 0x253;
        do {
          if (((pbVar11[iVar5 + 0x294] == 0x30) && (*pbVar9 == 0x39)) &&
             ((sVar4 = FUN_10017090((int)(pbVar11 + -0x95),(int)(pbVar9 + local_c + 1)), sVar4 == 1
              || (sVar4 == 2)))) {
            *pbVar9 = 0x4c;
          }
          iVar5 = iVar5 + 1;
          pbVar9 = pbVar9 + 1;
        } while (iVar5 < (int)(*pbVar11 - 1));
      }
      local_1c = local_1c + 1;
      local_c = local_c + -0x3c0;
      pbVar11 = pbVar11 + 0x3c0;
    } while (local_1c < *(short *)(iVar3 + 2));
  }
  iVar5 = 0;
  if (*(short *)(iVar3 + 2) != 1 && -1 < *(short *)(iVar3 + 2) + -1) {
    iVar10 = iVar3 + 0x934;
    do {
      if (*(char *)(iVar10 + 0xd5) == ']') {
        bVar1 = *(byte *)(iVar10 + -0x253);
        iVar12 = bVar1 - 1;
        sVar4 = FUN_10017090(iVar10 + 0xd8,0);
        if (-1 < sVar4) {
          uVar6 = FUN_10017030(iVar10 + -0x2e8,bVar1 - 2);
          if (-1 < (short)uVar6) {
            if (*(char *)(iVar10 + 0x3c0) != 'C') {
              if (*(char *)(iVar10 + iVar12) == '9') {
                *(undefined1 *)(iVar10 + iVar12) = 0x4d;
              }
              else if (*(char *)(iVar10 + iVar12) == '\x15') {
                *(undefined1 *)(iVar10 + iVar12) = 0x4e;
              }
            }
            if (*(char *)(iVar10 + iVar12) == '*') {
              *(undefined1 *)(iVar10 + iVar12) = 0x46;
            }
            else if (*(char *)(iVar10 + iVar12) == '5') {
              *(undefined1 *)(iVar10 + iVar12) = 0x47;
            }
          }
          if ((0 < iVar12) && ((&DAT_1007be3c)[*(char *)(iVar10 + -1 + iVar12)] == '\x01')) {
            cVar2 = *(char *)(iVar10 + iVar12);
            if (cVar2 == '9') {
              *(undefined1 *)(iVar10 + iVar12) = 0x48;
            }
            else if (cVar2 == '*') {
              *(undefined1 *)(iVar10 + iVar12) = 0x46;
            }
            else if (cVar2 == '5') {
              *(undefined1 *)(iVar10 + iVar12) = 0x47;
            }
          }
          if (((*(char *)(iVar10 + iVar12) == '\x15') || (*(char *)(iVar10 + iVar12) == '9')) &&
             (iVar7 = FUN_100560a0(&PTR_DAT_1007beec,*(char **)(iVar10 + 0x3b8),10,0x49),
             iVar7 != -1)) {
            if (*(char *)(iVar10 + iVar12) == '\x15') {
              *(undefined1 *)(iVar10 + iVar12) = 0x4f;
            }
            else if (*(char *)(iVar10 + iVar12) == '9') {
              *(undefined1 *)(iVar10 + iVar12) = 0x50;
            }
          }
        }
      }
      iVar5 = iVar5 + 1;
      iVar10 = iVar10 + 0x3c0;
    } while (iVar5 < *(short *)(iVar3 + 2) + -1);
  }
  return 1;
}

 (GhidraScript)  
INFO  DecompileNamedFunctions.java> ===== 0x100560a0 ===== (GhidraScript)  
INFO  DecompileNamedFunctions.java> Function: FUN_100560a0 @ 100560a0 (GhidraScript)  
INFO  DecompileNamedFunctions.java> 
int __cdecl FUN_100560a0(undefined4 *param_1,char *param_2,int param_3,short param_4)

{
  int iVar1;
  int iVar2;
  int iVar3;
  int iVar4;
  bool bVar5;
  
  if ((param_2 != (char *)0x0) && (*param_2 != '\0')) {
    iVar3 = 0;
    bVar5 = param_4 != 0x49;
    _param_4 = FUN_1001c2c0;
    iVar4 = param_3 + -1;
    if (bVar5) {
      _param_4 = _strcmp;
    }
    if (param_3 < 0) {
      if (-1 < iVar4) {
        do {
          iVar1 = (*_param_4)(*param_1,param_2);
          if (iVar1 == 0) {
            return iVar3;
          }
          iVar3 = iVar3 + 1;
          param_1 = param_1 + 1;
        } while (iVar3 <= iVar4);
      }
    }
    else {
      iVar1 = iVar4;
      if (-1 < iVar4) {
        do {
          iVar1 = iVar1 / 2;
          iVar2 = (*_param_4)(param_1[iVar1],param_2);
          if (iVar2 == 0) {
            return iVar1;
          }
          if (iVar2 < 0) {
            iVar3 = iVar1 + 1;
          }
          else {
            iVar4 = iVar1 + -1;
          }
          iVar1 = iVar4 + iVar3;
        } while (iVar3 <= iVar4);
        return -1;
      }
    }
  }
  return -1;
}

 (GhidraScript)  
INFO  REPORT: Save succeeded for processed file: /vt_pau.dll (HeadlessAnalyzer)  

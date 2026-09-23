#ifdef __cplusplus
extern "C"                          // assume C declarations for C++
{
#endif

extern short VT_Load_ENG(char *pszParam, HWND hWnd, int nSpeakerID, char *db_path, char *licensefile);
extern void  VT_UnLoad_ENG(int nSpeakerID);

#ifdef __cplusplus
}
#endif
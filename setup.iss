[Setup]
AppId={{38599DE5-0295-498E-94DA-EFC66F72F6EB}}
AppName=Blindfold Keeper
AppVersion=1.0
; C:\Program Files\Blindfold
DefaultDirName={pf}\Blindfold
OutputBaseFilename=Blindfold-Keeper
PrivilegesRequired=admin
SolidCompression=yes
WizardStyle=modern
OutputDir=output\windows

[Files]
Source: "win-blindfold-keeper.exe"; DestDir: "{app}"

[Registry]
Root: HKLM; Subkey: "SOFTWARE\Google\Chrome\NativeMessagingHosts\com.blindfold.keeper"; ValueType: string; ValueName: ""; ValueData: "{app}\com.blindfold.keeper.json"; Flags: uninsdeletekey

[Code]
// 동적으로 JSON 매니페스트 파일을 생성
procedure CurStepChanged(CurStep: TSetupStep);
var
  JsonPath: string;
  JsonContent: TStringList;
  AppPath: string;
begin
  if CurStep = ssPostInstall then
  begin
    // 실제 설치된 경로 (예: C:\Program Files\Blindfold)
    AppPath := ExpandConstant('{app}');
    
    // JSON의 "path" 값에 쓰기 위해 백슬래시(\)를 이스케이프 (\\)
    StringChange(AppPath, '\', '\\');

    // 생성할 JSON 파일의 전체 경로
    JsonPath := ExpandConstant('{app}\com.blindfold.keeper.json');

    // JSON 파일 내용 동적 생성
    JsonContent := TStringList.Create;
    try
      JsonContent.Add('{');
      JsonContent.Add('  "name": "com.blindfold.keeper",');
      JsonContent.Add('  "description": "Blindfold Device Key Storage",');
      JsonContent.Add('  "path": "' + AppPath + '\\win-blindfold-keeper.exe",');
      JsonContent.Add('  "type": "stdio",');
      JsonContent.Add('  "allowed_origins": [');
      JsonContent.Add('    "chrome-extension://cmgjlocmnppfpknaipdfodjhbplnhimk/"');
      JsonContent.Add('  ]');
      JsonContent.Add('}');
            
      JsonContent.SaveToFile(JsonPath);
    finally
      JsonContent.Free;
    end;
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
var
  JsonPath: string;
begin
  if CurUninstallStep = usUninstall then
  begin
    JsonPath := ExpandConstant('{app}\com.blindfold.keeper.json');
    DeleteFile(JsonPath);
  end;
end;
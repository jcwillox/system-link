directory = CreateObject("Scripting.FileSystemObject").GetParentFolderName(WScript.ScriptFullName)
command = directory & "\system-bridge.exe"

Set shell = CreateObject("WScript.Shell")
shell.CurrentDirectory = directory
shell.Run command, 0, False

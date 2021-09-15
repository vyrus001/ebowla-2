# Ebowla-2
reboot of https://github.com/Genetic-Malware/Ebowla in order to simplify / modernize the codebase and provide ongoing support

# Building / Running
the following is an example of compiling a hello world golang exe, packaging it into an encrypted package, and finally, building the implant with the package embedded within.
```
cd example
go build
mv example.exe ..\package\example.exe
cd ..\package\
go run .\main.go -s 'C:\Windows\System32\AboveLockAppHost.dll' -p .\example.exe
mv .\package ..\implant\
cd ..\implant\
go build -ldflags "-X main.seedPath=C:\Windows\System32"
.\implant.exe
```

### Concept Presentation Resources
Slides:

* Infiltrate 2016: https://github.com/Genetic-Malware/Ebowla/raw/master/Infiltrate_2016_Morrow_Pitts_Genetic_Malware.pdf
* Ekoparty 2016: https://github.com/Genetic-Malware/Ebowla/blob/master/Eko_2016_Morrow_Pitts_Master.pdf

Demos:
* https://www.youtube.com/watch?v=rRm3O7w5GHg
* https://youtu.be/Bu_qDrbX9Zo
* https://youtu.be/mlh70LtwmDo
* https://youtu.be/lyedtAtATGc (PowerShell)

# Payload Support
|Shared Object|
|:-----|
|Reflective DLL|

### Contributors
Original Project:
* https://github.com/wired33 (wrote most of the golang payload code)
* https://github.com/secretsquirrel (wrote the python payload code and most of the encryption code)

This Project:
* https://github.com/secretsquirrel

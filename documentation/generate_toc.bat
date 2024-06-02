@ECHO off
REM this uses Markdown TOC generator (to be installed separately) 
REM https://github.com/ycd/toc

REM all paths are relative to this directory. May be different on other machines

@ECHO ON
..\..\17.MD-TOC-generator\toc.exe -p ..\README.md
..\..\17.MD-TOC-generator\toc.exe -p ..\documentation\API.md
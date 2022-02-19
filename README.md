# svm
Yet another version manager, written in Go, to manage Apache Spark installations.

- Installs any version of Spark from https://archive.apache.org/dist/spark/
- Easy switching between different Spark versions

## Getting started
1. Download `svm` in a directory that is in `PATH`
2. Add the following lines to your `.zshrc`/`.bashrc` file:
```export PATH="$HOME/.svm/active/bin/:$PATH"```
3. Install spark version of your choice e.g. `svm install 2.2.2-with-hadoop-2.7`
4. Set spark version `svm use 2.2.2-with-hadoop-2.7`

## Example
```
❯ svm install 2.2.2-hadoop2.7
Fetching from: https://archive.apache.org/dist/spark/spark-2.2.2/spark-2.2.2-bin-hadoop2.7.tgz
191.44 MiB / 191.44 MiB 100 % [==============================================================] 0s ] 12.53 MiB/s
Installed 2.2.2-hadoop2.7 successfully
❯ svm use 2.2.2-hadoop2.7
Active Spark version: 2.2.2-hadoop2.7
❯ spark-shell --version
Welcome to                                 
      ____              __                 
     / __/__  ___ _____/ /__               
    _\ \/ _ \/ _ `/ __/  '_/               
   /___/ .__/\_,_/_/ /_/\_\   version 2.2.2
      /_/                                  
                                           
Using Scala version 2.11.8, OpenJDK 64-Bit Server VM, 11.0.11
Branch                                                       
Compiled by user  on 2018-06-27T14:30:46Z                    
Revision                                                     
Url                                                          
Type --help for more information.      
```
Done and ready to go!

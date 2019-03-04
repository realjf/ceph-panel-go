# ceph-panel-go
A Ceph Cluster Management Panel

### 获取librados依赖
#### librados for C/C++
Debian/Ubuntu

```sh
sudo apt-get install librados-dev
```
RHEL/CentOS
```sh
sudo yum install librados2-devel
```
检查是否安装成功
```sh
ls /usr/include/rados
```

#### librados for python
Debian/Ubuntu

```sh
sudo apt-get install python-rados
```
RHEL/CentOS
```sh
sudo yum install python-rados
```

#### librados for java
**1. 第一步：安装jna.jar包**

Debian/Ubuntu

```sh
sudo apt-get install libjna-java
```
RHEL/CentOS
```sh
sudo yum install jna
```
> jar文件位于/usr/share/java

**2. 第二步：克隆rados-java软件库**

```sh
git clone --recursive https://github.com/ceph/rados-java.git
```
**3. 第三步：构建rados-java软件库**
```sh
cd rados-java
ant
```
jar文件位于rados-java/target

** 4. 第四步：把rados的jar文件复制到统一位置（/usr/share/java），并确保它和jna jar都位于jvm路径里 **
```sh
sudo cp target/rados-0.1.3.jar /usr/share/java/rados-0.1.3.jar
sudo ln -s /usr/share/java/jna-3.2.7.jar /usr/lib/jvm/default-java/jre/lib/ext/jna-3.2.7.jar
sudo ln -s /usr/share/java/rados-0.1.3.jar  /usr/lib/jvm/default-java/jre/lib/ext/rados-0.1.3.jar
```

#### librados for php
**1. 安装php-dev**
Debian/Ubuntu
```sh
sudo apt-get install php5-dev build-essential
```
CentOS/RHEL
```sh
sudo yum install php-devel
```
**2. 克隆phprados源码库**
```sh
git clone https://github.com/ceph/phprados.git
```
**3. 构建phprados**
```sh
cd phprados
phpize
./configure
make
sudo make install
```
**4. 把下列配置加入php.ini以启用phprados**
```sh
extension=rados.so
```


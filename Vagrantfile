# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "centos"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # NOTE: This will enable public access to the opened port
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine and only allow access
  # via 127.0.0.1 to disable public access
  # config.vm.network "forwarded_port", guest: 80, host: 8080, host_ip: "127.0.0.1"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"
  config.vm.boot_timeout = 300

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"
  config.vm.synced_folder "/Users/lijinya/Yandex.Disk.localized/vagrant/", "/home/vagrant/data", create:true,
  :owner => "vagrant",
  :group => "vagrant",
  :mount_options => ["dmode=775","fmode=775"]


  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  config.vm.provider "virtualbox" do |vb|
    # Display the VirtualBox GUI when booting the machine
    #vb.gui = true
  
    # Customize the amount of memory on the VM:
    vb.memory = "2048"
  end

  #config.vm.provision :shell, :inline => "sudo sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/g' /etc/ssh/sshd_config; sudo systemctl restart sshd;", run: "always"
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  config.vm.provision "shell", inline: <<-SHELL
    #sudo yum install net-tools -y
    #sudo yum install expect -y
    #sudo yum install python-setuptools -y
    #sudo easy_install pip
    #sudo pip install pyotp
    #sudo pip install python-config
    #sudo yum install -y ncurses-devel lua-devel
    #sudo su - vagrant
    #sudo yum install git -y
    #cd data/
    #cp * ~vagrant/.ssh/
    #cd ~vagrant/.ssh/
    #chmod 600 config id_rsa
    #chmod 600 config id_rsa
    #chmod 644 id_rsa.pub known_hosts
    #cd ~
    #cp ~vagrant/data/jump.sh .
    #chmod +x jump.sh
    #chmod 600 config id_rsa
    #chmod 644 id_rsa.pub known_hosts
    #cd ~vagrant
    #sudo cp ~vagrant/data/jump.sh .
    #rsync -avz dev:/data/home/lijinya/.bashrc ~vagrant/
    #source ~vagrant/.bashrc
    #git clone https://github.com/magicmonty/bash-git-prompt.git ~/.bash-git-prompt
    #git clone git@github.com:vim/vim.git ~/vim
    #git clone git@github.com:ljy1814/vim.git ~/mvim
    #cd ~vagrant/mvim/
    #gcb golang origin/golang
    #ln -s `pwd`/vimrc2 ~/.vimrc
    #cd ~vagrant/vim
    #./configure --with-features=huge --enable-pythoninterp --enable-rubyinterp --enable-luainterp  --with-python-config-dir=/usr/lib/python2.7/config/ --enable-gui=gtk2 --enable-cscope --prefix=/usr
    #make
    #sudo make install
    #cd
    #cd ~vagrant/mvim/
    #sh install.sh
    #vim +PluginInstall +qall
    whoami
  SHELL
end


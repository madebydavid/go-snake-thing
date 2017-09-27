Vagrant.configure(2) do |config|

  config.vm.box = "debian/jessie64"
  config.vm.box_version = "8.9.0"

  config.vm.network "forwarded_port", guest: 8080, host: 8080

  # Increase the VM RAM from default 512MB to 2GB
  config.vm.provider :virtualbox do |vb|
    vb.memory = "2048"
  end

  # Base provisioning script
  config.vm.provision :shell,
    :keep_color => true,
    :path => "provisioning/base.sh"

  # Synced folder
  config.vm.synced_folder ".", "/vagrant",
    :type => "virtualbox", 
    :owner => "vagrant",
    :group => "vagrant"

end
# cuttle

Issuing certificates with fabric-ca

用fabric-ca代替cryptogen实现证书的颁发，现在完全替代cryptogen功能，后续会实现颁发特定证书

- 如果想颁发联盟中所有证书（即cryptogen所实现），需要配置static\crypto-config.yaml文件
该文件仿照fabric 中crypto-config.yaml文件，但不同的在于每个组织需要制定各自根ca的配置文件，并且在Subject中定义通用Subject属性

- 如果想颁发特定文件，需要配置static\cuttle.yaml文件


在main函数中定义好了RunConfig和RunSpeConfig两种实现方式
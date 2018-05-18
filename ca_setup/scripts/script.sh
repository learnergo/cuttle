#!/bin/bash
# Copyright London Stock Exchange Group All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
echo
echo " ____    _____      _      ____    _____  "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _| "
echo "\___ \    | |     / _ \   | |_) |   | |   "
echo " ___) |   | |    / ___ \  |  _ <    | |   "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|   "
echo

EnrollAdmin(){
	export FABRIC_CA_CLIENT_HOME=/etc/hyperledger/fabric-ca-client/Admins/ca.example.com
	fabric-ca-client enroll -u http://admin:adminpw@ca.example.com:7054
	mv /etc/hyperledger/fabric-ca-client/Admins/ca.example.com/msp/keystore/* /etc/hyperledger/fabric-ca-client/Admins/ca.example.com/msp/keystore/key.pem
	export FABRIC_CA_CLIENT_HOME=/etc/hyperledger/fabric-ca-client/Admins/tlsca.example.com
	fabric-ca-client enroll -u http://admin:adminpw@tlsca.example.com:7054
	mv /etc/hyperledger/fabric-ca-client/Admins/tlsca.example.com/msp/keystore/* /etc/hyperledger/fabric-ca-client/Admins/tlsca.example.com/msp/keystore/key.pem
	
	export FABRIC_CA_CLIENT_HOME=/etc/hyperledger/fabric-ca-client/Admins/ca.org1.example.com
	fabric-ca-client enroll -u http://admin:adminpw@ca.org1.example.com:7054
	mv /etc/hyperledger/fabric-ca-client/Admins/ca.org1.example.com/msp/keystore/* /etc/hyperledger/fabric-ca-client/Admins/ca.org1.example.com/msp/keystore/key.pem
	export FABRIC_CA_CLIENT_HOME=/etc/hyperledger/fabric-ca-client/Admins/tlsca.org1.example.com
	fabric-ca-client enroll -u http://admin:adminpw@tlsca.org1.example.com:7054
	mv /etc/hyperledger/fabric-ca-client/Admins/tlsca.org1.example.com/msp/keystore/* /etc/hyperledger/fabric-ca-client/Admins/tlsca.org1.example.com/msp/keystore/key.pem
	
	export FABRIC_CA_CLIENT_HOME=/etc/hyperledger/fabric-ca-client/Admins/ca.org2.example.com
	fabric-ca-client enroll -u http://admin:adminpw@ca.org2.example.com:7054
	mv /etc/hyperledger/fabric-ca-client/Admins/ca.org2.example.com/msp/keystore/* /etc/hyperledger/fabric-ca-client/Admins/ca.org2.example.com/msp/keystore/key.pem
	export FABRIC_CA_CLIENT_HOME=/etc/hyperledger/fabric-ca-client/Admins/tlsca.org2.example.com
	fabric-ca-client enroll -u http://admin:adminpw@tlsca.org2.example.com:7054
	mv /etc/hyperledger/fabric-ca-client/Admins/tlsca.org2.example.com/msp/keystore/* /etc/hyperledger/fabric-ca-client/Admins/tlsca.org2.example.com/msp/keystore/key.pem
	
	sleep 10
}

echo "Enroll Admin..."
EnrollAdmin

echo
echo " _____   _   _   ____   "
echo "| ____| | \ | | |  _ \  "
echo "|  _|   |  \| | | | | | "
echo "| |___  | |\  | | |_| | "
echo "|_____| |_| \_| |____/  "
echo

exit 0

echo
echo
echo '########## Cleaning everything (docker images, volumes, everything...) ##########'
sleep 2s
docker rm -vf $(docker ps -a -q) 
docker system prune
docker builder prune -a
docker volume prune --filter all=1
sleep 2s

echo
echo
echo '########## creating ca materials for Laptop 1 ##########'
sleep 2s
cd ./setup1/vm1
sudo rm -rf crypto-config
sleep 1s
cd ./create-certificate-with-ca
sudo rm -rf fabric-ca
docker-compose up -d
sudo chmod 777 *
sleep 1s
./create-certificate-with-ca.sh
cd ../../

echo
echo
echo '########## creating ca materials for Laptop 2 ##########'
sleep 2s
cd ./vm2
sudo rm -rf crypto-config
sleep 1s
cd ./create-certificate-with-ca
sudo rm -rf fabric-ca
docker-compose up -d
sudo chmod 777 *
sleep 1s
./create-certificate-with-ca.sh
cd ../../

echo
echo
echo '########## creating ca materials for Pi 4 ##########'
sleep 2s
cd ./vm3
sudo rm -rf crypto-config
sleep 1s
cd ./create-certificate-with-ca
sudo rm -rf fabric-ca
docker-compose up -d
sudo chmod 777 *
sleep 1s
./create-certificate-with-ca.sh
cd ../../

echo 
echo
echo '########## creating ca materials for Pi 5 ##########'
sleep 2s
cd ./vm4
sudo rm -rf crypto-config
sleep 1s
cd ./create-certificate-with-ca
sudo rm -rf fabric-ca
docker-compose up -d
sudo chmod 777 *
sleep 1s
./create-certificate-with-ca.sh
cd ../../

echo
echo
echo '########## show all running containers ##########'
docker ps

echo
echo
echo '########## bootstrap the blockchain, creating genesis block and anchor peers ##########'
cd ../artifacts/channel
sudo chmod 777 *
sudo rm genesis.block
sudo rm mychannel.tx
sudo rm Org1MSPanchors.tx
sudo rm Org2MSPanchors.tx
sudo rm Org3MSPanchors.tx
./create-artifacts.sh

# echo
# echo
# echo '########## running Laptop 1 nodes ##########'
# sleep 2s
# cd ../../setup1/vm1
# sudo chmod 777 *
# docker-compose up -d
# sleep 2s

# echo
# echo
# echo '########## running Laptop 2 nodes ##########'
# sleep 2s
# cd ../vm2
# sudo chmod 777 *
# docker-compose up -d
# sleep 2s

# echo
# echo
# echo '########## running Pi 4 nodes ##########'
# sleep 2s
# cd ../vm3
# sudo chmod 777 *
# docker-compose up -d
# sleep 2s

# echo
# echo
# echo '########## running Pi 5 nodes ##########'
# sleep 2s
# cd ../vm4
# sudo chmod 777 *
# docker-compose up -d
# sleep 2s

# echo
# echo
# echo '########## show all running containers ##########'
# docker ps

# echo
# echo
# echo '########## Laptop 1 nodes joining the channel ##########'
# sleep 1s
# cd ../vm1
# ./createChannel.sh

# echo
# echo
# echo '########## Laptop 2 nodes joining the channel ##########'
# sleep 1s
# cd ../vm2
# ./joinChannel.sh

# echo
# echo
# echo '########## Pi 4 nodes joining the channel ##########'
# sleep 1s
# cd ../vm3
# ./joinChannel.sh

# echo
# echo
# echo '########## Installing chaincode ##########'
# sleep 1s
# cd ../vm1/
# ./deployChaincode.sh

# cd ../vm2/
# ./installAndApproveChaincode.sh

# cd ../vm3/
# ./installAndApproveChaincode.sh


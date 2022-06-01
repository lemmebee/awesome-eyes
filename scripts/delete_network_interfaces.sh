#!/bin/bash
echo "Running delete_network_interfaces.sh ...."

stateFilePath=$1

aws eks delete-nodegroup \
--cluster-name `terraform output -raw -state $stateFilePath cluster_name` \
--nodegroup-name `terraform output -raw -state $stateFilePath node_group_name`

# sleep 200

eniAttachIds=$(aws ec2 describe-network-interfaces | grep -o '"AttachmentId": "[^"]*' | grep -o '[^"]*$')
for eniAttachId in $eniAttachIds; do echo "Deleting $eniAttachId..."; aws ec2 detach-network-interface --attachment-id $eniAttachId; done

for count in 1 2; do eniIds=$(aws ec2 describe-network-interfaces |grep -o '"NetworkInterfaceId": "[^"]*' | grep -o '[^"]*$'); for eniId in $eniIds; do echo "Deleting $eniId..."; aws ec2 delete-network-interface --network-interface-id $eniId; done; done
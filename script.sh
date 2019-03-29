#!/bin/sh

TERRAFORM_TFSTATE="terraform.tfstate"
TERRAFORM_TFSTATE_BACK="terraform.tfstate_back"
TFSTATE_FILE=$1/$TERRAFORM_TFSTATE

echo "Looking for terraform.tfstate file in $1"

cp $TFSTATE_FILE $1/$TERRAFORM_TFSTATE_BACK
echo "Created a back up of the state file, $TERRAFORM_TFSTATE_BACK in $1"

sed -i '' "/deployment_configuration.%/d" $TFSTATE_FILE
sed -i -e 's/deployment_configuration.//1' $TFSTATE_FILE
sed -i -e 's/catalog_configuration/deployment_configuration/1' $TFSTATE_FILE
sed -i -e 's/catalog_id/catalog_item_id/1' $TFSTATE_FILE
sed -i -e 's/catalog_name/catalog_item_name/1' $TFSTATE_FILE
sed -i -e 's/vra7_resource/vra7_deployment/1' $TFSTATE_FILE

echo "Successfully migrated the old state file to the new format"

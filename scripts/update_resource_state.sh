#!/bin/sh

TERRAFORM_TFSTATE="terraform.tfstate"
TERRAFORM_TFSTATE_BACK="terraform.tfstate_back"

echo "Looking for terraform.tfstate file"
if [ ! -f "$TERRAFORM_TFSTATE" ]
then
        echo "error, no terraform state file found"
        exit
fi

cp $TFSTATE_FILE $TERRAFORM_TFSTATE_BACK
echo "Created a back up of the state file, $TERRAFORM_TFSTATE_BACK"

if [ `uname` == 'Darwin' ]
then
    sed -i "" "/deployment_configuration.%/d" $TERRAFORM_TFSTATE
    sed -i "" 's/deployment_configuration.//1' $TERRAFORM_TFSTATE
    sed -i "" 's/catalog_configuration/deployment_configuration/1' $TERRAFORM_TFSTATE
    sed -i "" 's/catalog_id/catalog_item_id/1' $TERRAFORM_TFSTATE
    sed -i "" 's/catalog_name/catalog_item_name/1' $TERRAFORM_TFSTATE
    sed -i "" 's/vra7_resource/vra7_deployment/1' $TERRAFORM_TFSTATE
else
    sed -i "deployment_configuration.%/d" $TERRAFORM_TFSTATE
    sed -i 's/deployment_configuration.//1' $TERRAFORM_TFSTATE
    sed -i 's/catalog_configuration/deployment_configuration/1' $TERRAFORM_TFSTATE
    sed -i 's/catalog_id/catalog_item_id/1' $TERRAFORM_TFSTATE
    sed -i 's/catalog_name/catalog_item_name/1' $TERRAFORM_TFSTATE
    sed -i 's/vra7_resource/vra7_deployment/1' $TERRAFORM_TFSTATE
fi


echo "Successfully migrated the old state file to the new format"

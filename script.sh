#!/bin/bash


json=resultado.json

subdominio=$(jq -r '.[0].user_answer' $json)
ttl=$(jq -r '.[1].user_answer' $json)
zona_de_cadastro=$(jq -r '.[2].user_answer' $json)
grupo_de_recurso=$(jq -r '.[3].user_answer' $json)


#echo "$subdominio"
#echo "$ttl"
#echo "$zona_de_cadastro"
#echo "$grupo_de_recurso"


az network dns record-set cname list --resource-group production \
                                    -z makesystem.com.br \
                                    --output json | jq '.[].fqdn' | resultado=$(grep $subdominio.$zona_de_cadastro) 

if [ $resultado ]; then
    echo "Existe"
else
sleep 3
echo "Vamos criar o conjunto de registros"

az network dns record-set cname set-record --resource-group $grupo_de_recurso \
                                            --zone-name $zona_de_cadastro \
                                            --record-set-name $subdominio \
                                            --ttl $ttl \
                                            --cname ingress-$grupo_de_recurso.makesystem.com.br
fi


#echo "Vamos criar o conjunto de registro (: "

#az network dns record-set cname create --name $subdominio \
#                                    --resource-group $grupo_de_recurso \
#                                    --ttl $ttl \
#                                    --zone-name $zona_de_cadastro \
                                    

#sleep 5


./dynamodump/dynamodump.py:
	git clone git@github.com:bchew/dynamodump.git;\
	pushd dynamodump;\
	pip install -r requirements.txt;\
	pip install flake8;\
	popd
# restore-table(tableName) - restores table content from backup.
# 
# 
restore-table = python dynamodump/dynamodump.py --skipThroughputUpdate -m restore -r ${AWS_REGION} --dataOnly -s $(1)

backup-all: ./dynamodump/dynamodump.py
	python dynamodump/dynamodump.py -m backup  -r ${AWS_REGION} -s ${ADAPTIVE_CLIENT_ID}*

restore-all: ./dynamodump/dynamodump.py
	python dynamodump/dynamodump.py -m restore -r ${AWS_REGION} --dataOnly -s "*"
# ${ADAPTIVE_CLIENT_ID}*

restore-table-user-objective: ./dynamodump/dynamodump.py
	./rename-backup.sh
	python dynamodump/dynamodump.py --skipThroughputUpdate -m restore -r ${AWS_REGION} --dataOnly -s ${ADAPTIVE_CLIENT_ID}_user_objective

restore-table-strategy-objectives: ./dynamodump/dynamodump.py
	python dynamodump/dynamodump.py --skipThroughputUpdate -m restore -r ${AWS_REGION} --dataOnly -s ${ADAPTIVE_CLIENT_ID}_strategy_objectives

rename-resource-user-objective:
	cd terraform;\
	terraform state mv 'aws_dynamodb_table.user_objectives' 'aws_dynamodb_table.user_objective_dynamodb_table'

backup-all-zip: backup-all
	tar -cvz dump -f $(shell date -Idate)-dump-${ADAPTIVE_CLIENT_ID}.tar.gz

backup-all-zip-upload-to-s3: backup-all-zip
	export BUCKET="adaptive-dump-backups" ;\
	aws s3 cp $(shell date -Idate)-dump-${ADAPTIVE_CLIENT_ID}.tar.gz s3://$${BUCKET}/$(shell date -Idate)-dump-${ADAPTIVE_CLIENT_ID}.tar.gz

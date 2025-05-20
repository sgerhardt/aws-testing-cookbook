# generate mocks
generate:
	mockery --name=GetObjectAPI \
			--dir=s3 \
            --inpackage \
		 	--testonly

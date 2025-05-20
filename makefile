# generate mocks
generate:
	mockery --name=GetObjectAPI \
			--dir=s3 \
            --inpackage \
		 	--testonly
		  	--filename=getobjectapi_mock_test.go \

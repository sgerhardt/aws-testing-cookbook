# generate mocks
generate:
	mockery --name=GetObjectAPI --dir=s3 --output=s3/mocks --outpkg=mocks --filename=s3_get_object_api_mock.go
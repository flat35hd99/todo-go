FROM amazon/aws-lambda-go:latest

# Copy function code
COPY main ${LAMBDA_TASK_ROOT}

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "main" ]
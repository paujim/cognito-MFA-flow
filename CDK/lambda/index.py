import os
import logging
import json
import boto3

logger = logging.getLogger()
logger.setLevel(logging.INFO)

client = boto3.client('lambda')
fn_name = os.getenv("FUNCTION_NAME")


def handler(event, context):

    logger.info('## EVENT')
    logger.info(event)

    try:

        key = event["Records"][0]["s3"]["object"]["key"]
        bucket = event["Records"][0]["s3"]["bucket"]["name"]
        # version = event["Records"][0]["s3"]["object"]["versionId"]
        response = client.update_function_code(
            FunctionName=fn_name,
            S3Bucket=bucket,
            S3Key=key,
        )

    except Exception as e:
        logger.info(str(e))

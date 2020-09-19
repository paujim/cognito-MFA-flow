import os
import logging
import json
import boto3

logger = logging.getLogger()
logger.setLevel(logging.INFO)

client = boto3.client('lambda')
code_pipeline = boto3.client('codepipeline')


def put_job_success(job, message):
    logger.info('Putting job success')
    logger.info(message)
    code_pipeline.put_job_success_result(jobId=job)


def put_job_failure(job, message):
    logger.info('Putting job failure')
    logger.info(message)
    code_pipeline.put_job_failure_result(
        jobId=job,
        failureDetails={'message': message, 'type': 'JobFailed'},
    )


def continue_job_later(job, message):
    continuation_token = json.dumps({'previous_job_id': job})
    logger.info('Putting job continuation')
    logger.info(message)
    code_pipeline.put_job_success_result(
        jobId=job,
        continuationToken=continuation_token,
    )


def handler(event, context):
    job_id = "Unknoun"
    try:
        logger.info("Getting job details")
        job_id = event['CodePipeline.job']['id']
        job_data = event['CodePipeline.job']['data']
        logger.info("Getting user parameters")
        user_parameters = job_data['actionConfiguration']['configuration']['UserParameters']
        params = json.loads(user_parameters)
        bucket = params['sourceBucket']
        key = params['sourceKey']
        function_name = params["functionName"]

        logger.info("Updating lambda source")
        response = client.update_function_code(
            FunctionName=function_name,
            S3Bucket=bucket,
            S3Key=key,
        )
        put_job_success(job_id, 'Source Updated')

    except Exception as e:
        logger.info(str(e))
        put_job_failure(job_id, 'Function exception: ' + str(e))

    logger.info("Complete")
    return "Complete."

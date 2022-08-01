import google.oauth2.service_account as sa

def creds(SCOPES=None):
    kwargs = {}
    if SCOPES:
        kwargs['scopes'] = SCOPES

    return sa.Credentials.from_service_account_file(
        '.local/hsk-analysis-342a7eb65032.json', **kwargs)

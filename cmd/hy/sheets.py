# Built-in
import collections
import dicts

# Vendor
from googleapiclient.discovery import build

# Local
import ansi
import auth
import tts

SCOPES = [
    'https://www.googleapis.com/auth/drive',
    'https://www.googleapis.com/auth/drive.file',
    'https://www.googleapis.com/auth/spreadsheets',
]

DOCUMENT_ID = '1ZIY2MgZ-a4lPKpymfIbiS7eaO2tw3RZDBh861KlfRFs'

def _val(v):
    if 'effectiveValue' not in v:
        return None
    ev = v['effectiveValue']
    try:
        return ev['numberValue']
    except:
        return ev['stringValue']

def _unpackSheets(sheets):
    return [
        ([
            [
                _val(value)
                for value in row['values']
            ]
            for row in sheet['data'][0]['rowData']
        ], sheet['properties']['sheetId'])
        for sheet in sheets
    ]

Status = collections.namedtuple('Status', 'key tab color title')
statuses = [
    Status(key = 'n', color=ansi.rgb(255,109,103), title="Not known"         , tab="两字符词"),
    Status(key = 'j', color=ansi.rgb(200,160,  0), title="Just pronunciation", tab="我能发音的两字符词"),
    Status(key = 'k', color=ansi.rgb( 95,249,103), title="Known"             , tab="学过的两字符词"),
    Status(key = 'l', color=ansi.rgb(125,132,255), title="Long-term memory"  , tab="很好学过的两字符词"),
]
statusForKey = {s.key: s for s in statuses}

class Connection:
    def __init__(self):
        creds = auth.creds(SCOPES)
        service = build('sheets', 'v4', credentials=creds)
        self._sheets = service.spreadsheets()

    def get(self, ranges):
        result = self._sheets.get(
            spreadsheetId=DOCUMENT_ID,
            includeGridData=True,
            ranges=ranges,
        ).execute()
        return _unpackSheets(result['sheets'])

    def batchUpdate(self, body):
        return self._sheets.batchUpdate(
            spreadsheetId=DOCUMENT_ID,
            body=body,
        ).execute()

    def fetchGroups(self, cell, tones):
        (groups, ids) = zip(*[
            ((g[0][0] or '').split(), id)
            for (g, id) in self.get([
                f'{s.tab}!{cell}'
                for s in statuses
            ])
        ])

        words = [w for g in groups[:3] for w in g]

        # Update audio cache.
        ttsconn = tts.Connection()
        for word in words:
            for voice in tts.VOICES:
                ttsconn.fetchWord(word, voice)

        unknown_words = set(words) - set(dicts.cedict.keys())
        mismatched_words = {
            w: t
            for w in words
            for t in [dicts.cedict[w]]
            if t and (tones not in t)
        }
        if unknown_words or mismatched_words:
            if unknown_words:
                print('Words unknown to cedict:', unknown_words)
            if mismatched_words:
                print(f'''Words mismatched to tones {tones}: {
                    ', '.join(
                        f"{w} = {'/'.join(t)}(hsk={'/'.join(h.tones for t in dicts.hskdict[w].values() for h in t) or '∅'})"
                        for (w, t) in mismatched_words.items()
                    )
                }''')
            return None

        return (groups, ids)

    def save(self, ids, x, y, wordses):
        return self.batchUpdate(body={
            'requests': [
                {
                    'updateCells': {
                        'fields': "userEnteredValue",
                        'range': {
                            "sheetId": ids[i],
                            "startColumnIndex": x,
                            "startRowIndex": y,
                            "endColumnIndex": x + 1,
                            "endRowIndex": y + 1,
                        },
                        'rows': [{
                            'values': [{
                                'userEnteredValue': {
                                    'stringValue': ' '.join(words),
                                },
                            }],
                        }]
                    },
                }
                for (i, words) in enumerate(wordses)
            ],
        })

    def checkAllCells(conn):
        print('Checking all cells...')
        for (t1, col) in zip('1234', 'ABCD'):
            for (t2, row) in zip('12345', '12345'):
                conn.fetchGroups(f'{col}{row}', f'{t1}{t2}')

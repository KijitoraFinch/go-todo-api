import requests

class TodoAPI:
    BASE_URL = 'http://localhost:8080/api/v1'

    @staticmethod
    def get_todo(fname):
        payload = {
            'filename': fname,
            'completed_only': 'false',
            'uncompleted_only': 'false',
            'priority': 'A',
            'created_date': None,
            'due_date': '2042-01-01'
        }
        response = requests.get(f'{TodoAPI.BASE_URL}/filter', params=payload)
        return response.json()

    @staticmethod
    def post_todo():
        with open('test.txt', 'rb') as file:
            files = {'file': file}
            response = requests.post(f'{TodoAPI.BASE_URL}/upload', files=files)
        return response.json()

if __name__ == "__main__":
    api = TodoAPI
    response = api.post_todo()
    print(response)
    response = api.get_todo(response['filename'])
    print(response)


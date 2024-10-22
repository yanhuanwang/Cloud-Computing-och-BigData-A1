from flask import Flask, render_template

app = Flask(__name__, static_folder='.', template_folder='.')

@app.route('/')
def list_expenses():
    return render_template('listExpense.html')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)

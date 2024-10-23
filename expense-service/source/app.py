from flask import Flask, request, jsonify
import psycopg2
import logging
from psycopg2.extras import RealDictCursor
from datetime import datetime
import time
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

DB_CONNECTION_STRING = "host=postgres user=postgres password=mysecretpassword dbname=expensedb"

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Database connection
def get_db_connection():
    retries = 10
    while retries > 0:
        try:
            conn = psycopg2.connect(DB_CONNECTION_STRING)
            return conn
        except psycopg2.OperationalError as e:
            logger.warning(f"Unable to connect to the database, retrying... ({retries} retries left)")
            retries -= 1
            time.sleep(5)
    raise Exception("Unable to connect to the database after multiple attempts.")

@app.route('/readiness', methods=['GET'])
def readiness():
    try:
        conn = get_db_connection()
        conn.close()
        return "Ready", 200
    except Exception as e:
        return "Database not ready", 503

@app.route('/add-expense', methods=['POST'])
def add_expense():
    data = request.get_json()
    if not data or 'username' not in data or 'description' not in data or 'amount' not in data:
        return "Invalid request payload", 400

    username = data['username']
    description = data['description']
    amount = data['amount']
    date = datetime.now()

    try:
        conn = get_db_connection()
        with conn.cursor() as cur:
            cur.execute(
                """
                INSERT INTO expenses (username, description, amount, date)
                VALUES (%s, %s, %s, %s)
                """,
                (username, description, amount, date)
            )
            conn.commit()
        conn.close()
        return "Expense added successfully!", 200
    except Exception as e:
        logger.error(f"Error saving expense to database: {e}")
        return "Error saving expense to database", 500

@app.route('/get-expenses', methods=['GET'])
def get_expenses():
    username = request.args.get('username')
    if not username:
        return "Username is required", 400

    try:
        conn = get_db_connection()
        with conn.cursor(cursor_factory=RealDictCursor) as cur:
            cur.execute("SELECT id, username, description, amount, date FROM expenses WHERE username = %s", (username,))
            expenses = cur.fetchall()
        conn.close()
        return jsonify(expenses), 200
    except Exception as e:
        logger.error(f"Error fetching expenses: {e}")
        return "Error fetching expenses", 500

@app.route('/update-expense', methods=['PUT'])
def update_expense():
    data = request.get_json()
    if not data or 'id' not in data or 'username' not in data or 'description' not in data or 'amount' not in data:
        return "Invalid request payload", 400

    expense_id = data['id']
    username = data['username']
    description = data['description']
    amount = data['amount']
    date = datetime.now()

    try:
        conn = get_db_connection()
        with conn.cursor() as cur:
            cur.execute(
                """
                UPDATE expenses
                SET description = %s, amount = %s, date = %s
                WHERE id = %s AND username = %s
                """,
                (description, amount, date, expense_id, username)
            )
            conn.commit()
        conn.close()
        return "Expense updated successfully!", 200
    except Exception as e:
        logger.error(f"Error updating expense: {e}")
        return "Error updating expense", 500

@app.route('/delete-expense', methods=['DELETE'])
def delete_expense():
    data = request.get_json()
    if not data or 'id' not in data or 'username' not in data:
        return "Invalid request payload", 400

    expense_id = data['id']
    username = data['username']

    try:
        conn = get_db_connection()
        with conn.cursor() as cur:
            cur.execute("DELETE FROM expenses WHERE id = %s AND username = %s", (expense_id, username))
            conn.commit()
        conn.close()
        return "Expense deleted successfully!", 200
    except Exception as e:
        logger.error(f"Error deleting expense: {e}")
        return "Error deleting expense", 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081, debug=True)
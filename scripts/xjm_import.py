import csv
import io
import json
import subprocess

SQLSERVER_INSTANCE = r"localhost\SQLEXPRESS"
SQLSERVER_DB = "xjm"

POSTGRES_CONTAINER = "gamelink-postgres"
POSTGRES_DB = "gamelink"
POSTGRES_USER = "gamelink"
TARGET_SCHEMA = "xjm_import"


TABLE_DEFS = [
    (
        "BadReviews",
        [
            "Id",
            "HandlerId",
            "Content",
            "CreateTime",
            "RelatedOrderId",
            "CreateOrgId",
            "Status",
            "UpdateTime",
            "CreateUserId",
            "UpdateUserId",
            "CreateUser",
            "UpdateUser",
            "IsDelete",
            "ExtJson",
        ],
        "badreviews",
        [
            "id",
            "handlerid",
            "content",
            "createtime",
            "relatedorderid",
            "createorgid",
            "status",
            "updatetime",
            "createuserid",
            "updateuserid",
            "createuser",
            "updateuser",
            "isdelete",
            "extjson",
        ],
        """
        CREATE TABLE xjm_import.badreviews (
            id BIGINT,
            handlerid BIGINT,
            content TEXT,
            createtime TIMESTAMP,
            relatedorderid BIGINT,
            createorgid BIGINT,
            status TEXT,
            updatetime TIMESTAMP,
            createuserid BIGINT,
            updateuserid BIGINT,
            createuser TEXT,
            updateuser TEXT,
            isdelete BOOLEAN,
            extjson TEXT
        )
        """,
    ),
    (
        "BossDeposits",
        [
            "Id",
            "BossName",
            "RechargeTime",
            "RemainingAmount",
            "BDStatus",
            "CreateOrgId",
            "Status",
            "CreateTime",
            "UpdateTime",
            "CreateUserId",
            "UpdateUserId",
            "CreateUser",
            "UpdateUser",
            "IsDelete",
            "ExtJson",
            "PlatformUserId",
            "PlatformUserName",
        ],
        "bossdeposits",
        [
            "id",
            "bossname",
            "rechargetime",
            "remainingamount",
            "bdstatus",
            "createorgid",
            "status",
            "createtime",
            "updatetime",
            "createuserid",
            "updateuserid",
            "createuser",
            "updateuser",
            "isdelete",
            "extjson",
            "platformuserid",
            "platformusername",
        ],
        """
        CREATE TABLE xjm_import.bossdeposits (
            id BIGINT,
            bossname TEXT,
            rechargetime TIMESTAMP,
            remainingamount NUMERIC(18,2),
            bdstatus TEXT,
            createorgid BIGINT,
            status TEXT,
            createtime TIMESTAMP,
            updatetime TIMESTAMP,
            createuserid BIGINT,
            updateuserid BIGINT,
            createuser TEXT,
            updateuser TEXT,
            isdelete BOOLEAN,
            extjson TEXT,
            platformuserid TEXT,
            platformusername TEXT
        )
        """,
    ),
    (
        "ExpenseRecords",
        [
            "Id",
            "ExpenseTime",
            "Person",
            "Project",
            "Amount",
            "Account",
            "ReviewStatus",
            "CreateOrgId",
            "Status",
            "CreateTime",
            "UpdateTime",
            "CreateUserId",
            "UpdateUserId",
            "CreateUser",
            "UpdateUser",
            "IsDelete",
            "ExtJson",
        ],
        "expenserecords",
        [
            "id",
            "expensetime",
            "person",
            "project",
            "amount",
            "account",
            "reviewstatus",
            "createorgid",
            "status",
            "createtime",
            "updatetime",
            "createuserid",
            "updateuserid",
            "createuser",
            "updateuser",
            "isdelete",
            "extjson",
        ],
        """
        CREATE TABLE xjm_import.expenserecords (
            id BIGINT,
            expensetime TIMESTAMP,
            person TEXT,
            project TEXT,
            amount NUMERIC(18,2),
            account TEXT,
            reviewstatus TEXT,
            createorgid BIGINT,
            status TEXT,
            createtime TIMESTAMP,
            updatetime TIMESTAMP,
            createuserid BIGINT,
            updateuserid BIGINT,
            createuser TEXT,
            updateuser TEXT,
            isdelete BOOLEAN,
            extjson TEXT
        )
        """,
    ),
    (
        "Orders",
        [
            "Id",
            "OStatus",
            "Remake",
            "ProjectName",
            "IsDeposits",
            "OrderAmount",
            "BossName",
            "Handler1",
            "Handler2",
            "CreateOrgId",
            "Status",
            "CreateTime",
            "UpdateTime",
            "CreateUserId",
            "UpdateUserId",
            "CreateUser",
            "UpdateUser",
            "IsDelete",
            "ExtJson",
            "OutsideOrderId",
            "SpareFlag",
        ],
        "orders",
        [
            "id",
            "ostatus",
            "remake",
            "projectname",
            "isdeposits",
            "orderamount",
            "bossname",
            "handler1",
            "handler2",
            "createorgid",
            "status",
            "createtime",
            "updatetime",
            "createuserid",
            "updateuserid",
            "createuser",
            "updateuser",
            "isdelete",
            "extjson",
            "outsideorderid",
            "spareflag",
        ],
        """
        CREATE TABLE xjm_import.orders (
            id BIGINT,
            ostatus TEXT,
            remake TEXT,
            projectname TEXT,
            isdeposits BOOLEAN,
            orderamount NUMERIC(18,2),
            bossname TEXT,
            handler1 BIGINT,
            handler2 BIGINT,
            createorgid BIGINT,
            status TEXT,
            createtime TIMESTAMP,
            updatetime TIMESTAMP,
            createuserid BIGINT,
            updateuserid BIGINT,
            createuser TEXT,
            updateuser TEXT,
            isdelete BOOLEAN,
            extjson TEXT,
            outsideorderid TEXT,
            spareflag TEXT
        )
        """,
    ),
    (
        "OrderSettlements",
        [
            "Id",
            "StartTime",
            "OrderId",
            "HandlerId",
            "HandlerAmount",
            "HandlerGameName",
            "PicUrl",
            "OSStatus",
            "Remake",
            "Commission",
            "EndTime",
            "CreateOrgId",
            "Status",
            "CreateTime",
            "UpdateTime",
            "CreateUserId",
            "UpdateUserId",
            "CreateUser",
            "UpdateUser",
            "IsDelete",
            "ExtJson",
        ],
        "ordersettlements",
        [
            "id",
            "starttime",
            "orderid",
            "handlerid",
            "handleramount",
            "handlergamename",
            "picurl",
            "osstatus",
            "remake",
            "commission",
            "endtime",
            "createorgid",
            "status",
            "createtime",
            "updatetime",
            "createuserid",
            "updateuserid",
            "createuser",
            "updateuser",
            "isdelete",
            "extjson",
        ],
        """
        CREATE TABLE xjm_import.ordersettlements (
            id BIGINT,
            starttime TIMESTAMP,
            orderid BIGINT,
            handlerid BIGINT,
            handleramount NUMERIC(18,2),
            handlergamename TEXT,
            picurl TEXT,
            osstatus TEXT,
            remake TEXT,
            commission NUMERIC(18,2),
            endtime TIMESTAMP,
            createorgid BIGINT,
            status TEXT,
            createtime TIMESTAMP,
            updatetime TIMESTAMP,
            createuserid BIGINT,
            updateuserid BIGINT,
            createuser TEXT,
            updateuser TEXT,
            isdelete BOOLEAN,
            extjson TEXT
        )
        """,
    ),
    (
        "PaySlips",
        [
            "PayslipId",
            "Id",
            "HandlerId",
            "Month",
            "SalaryDetails",
            "Commission",
            "IsWithdrawn",
            "CreateOrgId",
            "Status",
            "CreateTime",
            "UpdateTime",
            "CreateUserId",
            "UpdateUserId",
            "CreateUser",
            "UpdateUser",
            "IsDelete",
            "ExtJson",
        ],
        "payslips",
        [
            "payslipid",
            "id",
            "handlerid",
            "month",
            "salarydetails",
            "commission",
            "iswithdrawn",
            "createorgid",
            "status",
            "createtime",
            "updatetime",
            "createuserid",
            "updateuserid",
            "createuser",
            "updateuser",
            "isdelete",
            "extjson",
        ],
        """
        CREATE TABLE xjm_import.payslips (
            payslipid INTEGER,
            id BIGINT,
            handlerid INTEGER,
            month TEXT,
            salarydetails TEXT,
            commission NUMERIC(18,2),
            iswithdrawn BOOLEAN,
            createorgid BIGINT,
            status TEXT,
            createtime TIMESTAMP,
            updatetime TIMESTAMP,
            createuserid BIGINT,
            updateuserid BIGINT,
            createuser TEXT,
            updateuser TEXT,
            isdelete BOOLEAN,
            extjson TEXT
        )
        """,
    ),
    (
        "RechargeRecord",
        [
            "Id",
            "BossName",
            "RechargeNum",
            "IsGift",
            "PicUrl",
            "RechargeTime",
            "Reporter",
            "Remark",
            "CreateOrgId",
            "Status",
            "CreateTime",
            "UpdateTime",
            "CreateUserId",
            "UpdateUserId",
            "CreateUser",
            "UpdateUser",
            "IsDelete",
            "ExtJson",
        ],
        "rechargerecord",
        [
            "id",
            "bossname",
            "rechargenum",
            "isgift",
            "picurl",
            "rechargetime",
            "reporter",
            "remark",
            "createorgid",
            "status",
            "createtime",
            "updatetime",
            "createuserid",
            "updateuserid",
            "createuser",
            "updateuser",
            "isdelete",
            "extjson",
        ],
        """
        CREATE TABLE xjm_import.rechargerecord (
            id BIGINT,
            bossname TEXT,
            rechargenum NUMERIC(18,2),
            isgift BOOLEAN,
            picurl TEXT,
            rechargetime TIMESTAMP,
            reporter TEXT,
            remark TEXT,
            createorgid BIGINT,
            status TEXT,
            createtime TIMESTAMP,
            updatetime TIMESTAMP,
            createuserid BIGINT,
            updateuserid BIGINT,
            createuser TEXT,
            updateuser TEXT,
            isdelete BOOLEAN,
            extjson TEXT
        )
        """,
    ),
    (
        "TheProjectDisk",
        [
            "Id",
            "ProjectType",
            "ProjectName",
            "ProjectAmount",
            "ProjectPictureUrl",
            "ProjectText",
            "SecretAmount",
            "SingleAmount",
            "CreateOrgId",
            "Status",
            "CreateTime",
            "UpdateTime",
            "CreateUserId",
            "UpdateUserId",
            "CreateUser",
            "UpdateUser",
            "IsDelete",
            "ExtJson",
        ],
        "theprojectdisk",
        [
            "id",
            "projecttype",
            "projectname",
            "projectamount",
            "projectpictureurl",
            "projecttext",
            "secretamount",
            "singleamount",
            "createorgid",
            "status",
            "createtime",
            "updatetime",
            "createuserid",
            "updateuserid",
            "createuser",
            "updateuser",
            "isdelete",
            "extjson",
        ],
        """
        CREATE TABLE xjm_import.theprojectdisk (
            id BIGINT,
            projecttype TEXT,
            projectname TEXT,
            projectamount NUMERIC(18,2),
            projectpictureurl TEXT,
            projecttext TEXT,
            secretamount NUMERIC(18,2),
            singleamount NUMERIC(18,2),
            createorgid BIGINT,
            status TEXT,
            createtime TIMESTAMP,
            updatetime TIMESTAMP,
            createuserid BIGINT,
            updateuserid BIGINT,
            createuser TEXT,
            updateuser TEXT,
            isdelete BOOLEAN,
            extjson TEXT
        )
        """,
    ),
]


def run_psql(sql: str) -> None:
    subprocess.check_call(
        [
            "docker",
            "exec",
            "-i",
            POSTGRES_CONTAINER,
            "psql",
            "-U",
            POSTGRES_USER,
            "-d",
            POSTGRES_DB,
            "-v",
            "ON_ERROR_STOP=1",
            "-c",
            sql,
        ]
    )


def copy_table(source_table: str, source_cols: list[str], target_table: str, target_cols: list[str]) -> None:
    select_sql = (
        "SELECT "
        + ", ".join(source_cols)
        + " FROM "
        + source_table
        + " FOR JSON PATH, INCLUDE_NULL_VALUES"
    )
    sqlcmd = [
        "sqlcmd",
        "-S",
        SQLSERVER_INSTANCE,
        "-d",
        SQLSERVER_DB,
        "-E",
        "-h",
        "-1",
        "-y",
        "8000",
        "-w",
        "32767",
        "-f",
        "65001",
        "-Q",
        f"SET NOCOUNT ON; {select_sql}",
    ]
    raw = subprocess.check_output(sqlcmd)
    if raw.startswith(b"\xef\xbb\xbf"):
        raw = raw[3:]
    json_text = raw.decode("utf-8", errors="replace").strip()
    if not json_text:
        rows = []
    else:
        json_text = "".join(line.strip() for line in json_text.splitlines() if line.strip())
        rows = json.loads(json_text)

    buffer = io.StringIO()
    writer = csv.writer(
        buffer,
        delimiter="|",
        quotechar='"',
        quoting=csv.QUOTE_MINIMAL,
        lineterminator="\n",
    )
    null_token = "__NULL__"
    for row in rows:
        writer.writerow([null_token if row.get(col) is None else row.get(col) for col in source_cols])
    data = buffer.getvalue().encode("utf-8")

    copy_sql = (
        f"COPY {TARGET_SCHEMA}.{target_table} ({', '.join(target_cols)}) "
        "FROM STDIN WITH (FORMAT csv, DELIMITER '|', QUOTE '\"', ESCAPE '\"', NULL '__NULL__')"
    )
    psql_cmd = [
        "docker",
        "exec",
        "-i",
        POSTGRES_CONTAINER,
        "psql",
        "-U",
        POSTGRES_USER,
        "-d",
        POSTGRES_DB,
        "-c",
        copy_sql,
    ]
    proc = subprocess.Popen(psql_cmd, stdin=subprocess.PIPE)
    proc.communicate(data)
    if proc.returncode != 0:
        raise SystemExit(proc.returncode)


def main() -> None:
    run_psql(f"CREATE SCHEMA IF NOT EXISTS {TARGET_SCHEMA}")
    for _, _, _, _, ddl in TABLE_DEFS:
        run_psql(f"DROP TABLE IF EXISTS {TARGET_SCHEMA}.{ddl.split()[2].split('.')[-1]}")
        run_psql(ddl)

    for source_table, source_cols, target_table, target_cols, _ in TABLE_DEFS:
        run_psql(f"TRUNCATE {TARGET_SCHEMA}.{target_table}")
        copy_table(source_table, source_cols, target_table, target_cols)


if __name__ == "__main__":
    main()

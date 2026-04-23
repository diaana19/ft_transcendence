# Bug Reporting Guidelines

To ensure efficient debugging and resolution, all reported issues must follow the format below.

---

## Required Information

### 1. Context

* Area: (Backend / Frontend)
* File:
* Function / Module:
* Branch (if applicable):

---

### 2. Description of the Issue

* Clear description of the problem:
* Expected behavior:
* Observed behavior:

---

### 3. Steps to Reproduce

Provide a precise sequence of steps:
1.
2.
3.

---

### 4. Test Data

* Input used:
* Command executed (backend):
* URL / Route (frontend):

---

### 5. Error Output

Provide the complete error message:

```id="j4v0u9"
(paste the full error output here)
```

---

### 6. Logs (Mandatory)

You must provide a compressed archive of the logs directory.

* Archive format: `.zip`, `.tar`, or `.tar.gz`
* The archive must contain the full `logs/` directory
* Do not modify or filter the log files

Example:

```id="y6c2wd"
logs.zip
└── logs/
    ├── backend.log
    ├── frontend.log
    └── error.log
```

---

### 7. Additional Information (Optional)

* Screenshots
* Network logs (frontend)
* Any relevant details

---

## Notes

* Reports that do not follow this format may be ignored
* Missing logs will result in the report being rejected
* Always specify whether the issue is backend or frontend
* Ensure the issue is reproducible before reporting
* Check that the issue has not already been reported
* Be precise and concise in your description

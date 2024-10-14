Here’s an extended and formatted version of your README.md for the RansomGuard project:

---

# RansomGuard

**RansomGuard** is a standalone ransomware detection and prevention solution designed for Linux servers. It leverages multiple detection methods and provides a rollback mechanism to recover potentially compromised files. 

## Key Features

RansomGuard employs **three simultaneous detection methods** to safeguard your system against ransomware attacks:

### 1. Trapping (Honeypot Strategy)
RansomGuard generates decoy files in critical directories of the system (e.g., `/usr`, `/etc`, etc.) with interesting file extensions such as `.doc`, `.pdf`, and `.db`. These files act as honeypots to trap potential ransomware activity. 

- The system continuously watches these honeypot files for any unauthorized modifications.
- Upon detecting changes, RansomGuard triggers an alert and takes immediate action to stop the malicious process.

### 2. Behavior Analysis (Kernel API Monitoring)
RansomGuard monitors the behavior of each active process by tracking their **Kernel API calls**. It looks for suspicious patterns, such as attempts to rename files, delete shadow copies, or encrypt files—common techniques used by ransomware.

- If any process exhibits unsecure behavior based on predefined API call chains, RansomGuard stops the process and sends an alert.

### 3. Static Analysis (Executable File Monitoring)
RansomGuard inspects each new static executable file that enters the system to detect suspicious characteristics based on a set of predefined keywords or patterns.

- If a potentially harmful executable file is detected, it is quarantined, and an alert is sent to the system administrator.

## Rollback Mechanism

In cases where Behavior and Static Analysis fail to detect a ransomware attack, but the Trapping method detects file tampering, RansomGuard provides a **rollback mechanism** to restore damaged files:

- **Rollback Configuration**: RansomGuard supports rollback using file system methods such as `btrfs` snapshots. Customers can configure their rollback strategy in the configuration file.
- **Automatic Restoration**: If a ransomware attack is confirmed through honeypot tampering, and rollback is enabled, RansomGuard will automatically restore the affected files from the specified backup or snapshot.

## Configuration

All configurations for RansomGuard are stored in a simple raw file (e.g., `.txt`), and all data is encrypted using a hardcoded key to ensure security. The configuration file includes:

- **Honeypot File Configuration**: File names, extensions, volume, and paths to generate honeypot files.
- **Email Addresses**: For sending alerts to the system administrator and RansomGuard company.
- **Keywords for Static Analysis**: A list of suspicious keywords for detecting malicious executable files.
- **Kernel API Monitoring**: Secure and unsecure API call chains for Behavior Analysis.
- **Rollback Configuration**: Rollback preferences, including whether rollback is enabled and the method (e.g., `btrfs` snapshots).

## Alert System

RansomGuard alerts are sent via email to two recipients:
1. The system administrator.
2. The RansomGuard company for additional support.

The alert system is triggered in the following situations:
- A honeypot file is modified.
- A suspicious Kernel API call chain is detected.
- A malicious executable file is detected through Static Analysis.
- The configuration file is modified.

## How to Use

1. **Download and Install**: Clone the repository and build the executable.
2. **Configure**: Modify the configuration file according to your server’s needs, specifying honeypot locations, alert emails, and rollback settings.
3. **Run RansomGuard**: Execute the standalone binary on your server. It will immediately start monitoring your system using all three detection methods.

## Roadmap

- Future versions may include support for additional rollback mechanisms and detection strategies.
- Further improvements in real-time Kernel monitoring and file analysis.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

### Example Configuration (config.json)
```json
{
  "honeypot_files": [
    {
      "name": "decoy",
      "extensions": [".doc", ".pdf", ".db"],
      "volume": "1MB",
      "route": "/usr"
    }
  ],
  "customer_email": "admin@example.com",
  "company_email": "ransomguard@example.com",
  "rollbackwant": true,
  "rollback_method": "btrfs",
  "rollback_variables": {
    "subvolume": "/backup/subvol",
    "snapshot": "/backup/snapshot"
  },
  "static_file_keywords": ["encrypt", "ransom", "lock"],
  "secure_processes": ["backup.sh", "rsync"],
  "unsecure_api_calls": ["rename", "unlink", "chmod"]
}
```

---

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any bugs, feature requests, or improvements.

---

This README should give potential users a good overview of what **RansomGuard** does, how it works, and how to configure it on their systems. Let me know if you need any additional changes!

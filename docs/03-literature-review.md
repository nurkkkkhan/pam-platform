# 3. Literature Review

## 3.1 Method

Sources were searched in IEEE Xplore, ACM Digital Library, ScienceDirect, arXiv and NIST CSRC using the keywords *privileged access management*, *least privilege*, *zero trust architecture*, *insider threat detection*, *secret management*, *tamper-evident logging*, *just-in-time access*. Inclusion criteria: (1) directly relevant to controlling, monitoring or auditing privileged access; (2) peer-reviewed publication, official standard, or widely cited industry report; (3) for the rapidly evolving topics (zero trust, insider-threat detection, secrets, market), published 2020 or later. Foundational works (least privilege, RBAC, secure logging) are kept regardless of year. 18 sources are used; 13 are scientific publications (peer-reviewed or arXiv preprints), 10 of the 18 are from 2021–2025.

## 3.2 Thematic review

### Theme 1 — Least privilege and access-control models

Saltzer and Schroeder formulated the *principle of least privilege*: every program and user should operate with the minimum privilege necessary [1]. Privileged accounts in practice violate this principle because they are permanent and shared. Role-based access control (RBAC) [2] provides the administrative model for assigning permissions through roles, and is still the baseline in all PAM products. However, RBAC alone is static; NIST's ABAC guide [18] shows how attributes (time, device, risk) refine decisions. **Implication for the project:** PAM Platform uses RBAC for *eligibility* (who *may* request access to which targets) and time-bound grants + contextual conditions (time window, MFA, approval) for *activation*, approximating ABAC without a full policy language.

### Theme 2 — Zero trust as the architectural frame for PAM

NIST SP 800-207 defines zero trust architecture (ZTA) with a policy decision point (PDP) and policy enforcement point (PEP) that authorize each session rather than trusting network location [5]. Syed et al. survey ZTA and highlight identity, continuous authentication and micro-segmentation as core building blocks, as well as open problems in performance and legacy integration [6]. Buck et al.'s multivocal review finds a gap between practitioner claims and academic evidence and calls for more empirical evaluation [7]. Teerakanok et al. focus on migration and identify the difficulty of moving existing infrastructure gradually [8]. Gambo and Almulhem's 2025 systematic review of 91 studies confirms a sharp growth of ZTA research since 2020 and repeatedly notes that many frameworks were validated only in simulation [9]. **Implication:** the platform's gateway is designed explicitly as a PEP and the control plane as a PDP; evaluation is done on a real test bench, not only in simulation, to answer the critique of [7] and [9]. Migration [8] is supported by allowing targets to be onboarded one by one without changing network topology.

### Theme 3 — Insider threat and monitoring of privileged sessions

Privileged users are a classic insider-threat population. Singh and Sharma's systematic review identifies privilege exploitation, data theft and sabotage as main insider attack types, and lists indicators such as off-hours access and unexpected privilege changes [13]. Yuan and Wu review deep-learning detectors and conclude that labelled data is scarce and most work relies on the synthetic CERT dataset [10], whose generation method is described by Glasser and Lindauer [17]. Le and Zincir-Heywood show that unsupervised ensembles can rank anomalous users without labels, but still produce false positives that require analyst review [11]. **Implication:** ML-based detection is not reliable enough to *block* sessions automatically; the MVP uses deterministic, explainable rules (blocked command patterns, off-hours, unusual target) for blocking, and records structured session data so ML detection can be added later (future work). This is a limitation the project consciously accepts.

### Theme 4 — Secret management

Basak et al. identified 24 secret-management practices from grey literature and found developers frequently hard-code or share secrets due to the lack of usable tooling [12]. Verizon's 2025 DBIR reports credential abuse as the top initial access vector (22%) [14]. NIST SP 800-53 requires authenticator management, account management and least privilege (IA-5, AC-2, AC-6) [15]. **Implication:** the vault must (a) never reveal credentials to users during brokered sessions, (b) rotate automatically after use, and (c) expose an API so applications stop hard-coding secrets.

### Theme 5 — Tamper-evident audit logs

A PAM system's audit log is only useful if administrators — the very users being audited — cannot silently alter it. Schneier and Kelsey proposed forward-secure hash chains for audit logs on untrusted machines [3]. Crosby and Wallach introduced history trees giving efficient proofs of log consistency [4]. NIST SP 800-53 control AU-9 requires protection of audit information [15]. **Implication:** PAM Platform implements a hash chain (each record contains `SHA-256(prev_hash || record)`) with periodic anchoring of the head hash outside the database; Merkle/history trees [4] are noted as a scalability improvement.

### Theme 6 — Market and existing solutions

The October 2025 Gartner Magic Quadrant for PAM names CyberArk, BeyondTrust and Delinea among the Leaders [16]. CyberArk was acquired by Palo Alto Networks (completed in 2026), which reinforces the trend of PAM becoming part of larger identity-security platforms [16]. Open-source tools (Teleport Community, JumpServer, Apache Guacamole, HashiCorp Boundary) cover parts of the problem; their licensing and feature limits are analysed in `04-existing-systems-comparison.md`.

## 3.3 Structured analysis table

| # | Source | Year | Type | Problem addressed | Approach / contribution | Limitations noted | Used in PAM Platform for |
|---|---|---|---|---|---|---|---|
| 1 | Saltzer & Schroeder | 1975 | Journal (Proc. IEEE) | Design principles for protection | Least privilege, complete mediation, economy of mechanism | Conceptual | Core design principles; complete mediation via gateway |
| 2 | Sandhu et al. | 1996 | Journal (IEEE Computer) | Scalable permission administration | RBAC model family (RBAC0–3) | Static, no context | Eligibility model |
| 3 | Schneier & Kelsey | 1999 | Journal (ACM TISSEC) | Logs on untrusted hosts | Hash-chained, forward-secure logs | Needs trusted verifier | Audit hash chain |
| 4 | Crosby & Wallach | 2009 | Conference (USENIX Security) | Efficient tamper-evident logging | History trees, membership & consistency proofs | Implementation complexity | Future scalability of audit |
| 5 | Rose et al. (NIST SP 800-207) | 2020 | Standard | Zero trust architecture | PDP/PEP logical model, deployment models | High-level, no reference implementation | Architecture (gateway = PEP) |
| 6 | Syed et al. | 2022 | Journal (IEEE Access) | ZTA survey | Taxonomy of ZTA components; authentication, access control | Performance & legacy integration open | Requirements on continuous auth, MFA |
| 7 | Buck et al. | 2021 | Journal (Computers & Security) | Knowledge gaps in zero trust | Multivocal literature review | Little empirical evidence | Motivation for real test-bench evaluation |
| 8 | Teerakanok et al. | 2021 | Journal (Security & Comm. Networks) | Migration to ZTA | Review of migration challenges | Organizational focus | Incremental onboarding of targets |
| 9 | Gambo & Almulhem | 2025 | arXiv preprint (SLR) | State of ZTA research 2016–2025 | PRISMA review of 91 studies, taxonomy | Many works simulation-only; preprint | Up-to-date evidence; evaluation design |
| 10 | Yuan & Wu | 2021 | Journal (Computers & Security) | Deep learning for insider threat | Review of DL detectors and datasets | Scarce real labelled data | Decision to use rules, not ML, for blocking |
| 11 | Le & Zincir-Heywood | 2021 | Journal (IEEE TNSM) | Unsupervised insider detection | Ensembles of unsupervised detectors | False positives need analyst | Future work: anomaly ranking |
| 12 | Basak et al. | 2022 | Conference (IEEE SecDev) | Secret management practices | Grey-literature review, 24 practices | Grey literature basis | Vault + secrets API requirements |
| 13 | Singh & Sharma | 2022 | arXiv preprint (SLR) | Insider threats | Review of 37 primary studies; threat types and countermeasures | Preprint; inconsistencies in counts | Detection rules (off-hours, privilege change) |
| 14 | Verizon DBIR | 2025 | Industry report | Breach statistics | Analysis of real incidents | Vendor dataset bias | Evidence of relevance |
| 15 | NIST SP 800-53 Rev.5 | 2020 | Standard | Security controls catalogue | AC-2, AC-6, AU-9, IA-5 controls | Generic, not product-specific | Compliance mapping, NFRs |
| 16 | Gartner MQ for PAM (via vendor releases); Palo Alto–CyberArk news | 2025–2026 | Industry analysis | PAM market | Vendor positioning | Paywalled; vendor-reported | Selection of systems to compare |
| 17 | Glasser & Lindauer | 2013 | Workshop (IEEE S&P Workshops) | Insider threat test data | Synthetic CERT dataset generation | Synthetic data | Understanding evaluation data limits |
| 18 | Hu et al. (NIST SP 800-162) | 2014 | Standard | ABAC | Attribute-based policy model | Complex to administer | Contextual conditions on grants |

## 3.4 What is known, and what this project takes into account

- **Known:** least privilege + RBAC are the foundation [1][2]; PAM fits the PEP/PDP model of zero trust [5][6]; credential misuse is the leading attack vector [14]; hash chains give practical tamper evidence [3].
- **Limitations accepted:** ML insider detection is not reliable enough to block sessions [10][11] → rule-based blocking; ZTA studies often lack real-world evaluation [7][9] → test-bench evaluation with a baseline; full ABAC is hard to administer [18] → limited contextual conditions.
- **Gap addressed:** a self-hosted, open, SME-scale platform that combines vault + rotation, brokered SSH/PostgreSQL access, JIT approval, session recording and verifiable audit in one product — currently split across several open-source tools or locked in enterprise licences (see §4).

## 3.5 References

[1] J. H. Saltzer and M. D. Schroeder, "The protection of information in computer systems," *Proceedings of the IEEE*, vol. 63, no. 9, pp. 1278–1308, 1975. doi:10.1109/PROC.1975.9939

[2] R. S. Sandhu, E. J. Coyne, H. L. Feinstein, and C. E. Youman, "Role-based access control models," *IEEE Computer*, vol. 29, no. 2, pp. 38–47, 1996. doi:10.1109/2.485845

[3] B. Schneier and J. Kelsey, "Secure audit logs to support computer forensics," *ACM Transactions on Information and System Security*, vol. 2, no. 2, pp. 159–176, 1999. doi:10.1145/317087.317089

[4] S. A. Crosby and D. S. Wallach, "Efficient data structures for tamper-evident logging," in *Proc. 18th USENIX Security Symposium*, 2009, pp. 317–334.

[5] S. Rose, O. Borchert, S. Mitchell, and S. Connelly, "Zero Trust Architecture," NIST Special Publication 800-207, 2020. doi:10.6028/NIST.SP.800-207

[6] N. F. Syed, S. W. Shah, A. Shaghaghi, A. Anwar, Z. Baig, and R. Doss, "Zero Trust Architecture (ZTA): A comprehensive survey," *IEEE Access*, vol. 10, pp. 57143–57179, 2022. doi:10.1109/ACCESS.2022.3174679

[7] C. Buck, C. Olenberger, A. Schweizer, F. Völter, and T. Eymann, "Never trust, always verify: A multivocal literature review on current knowledge and research gaps of zero-trust," *Computers & Security*, vol. 110, 102436, 2021. doi:10.1016/j.cose.2021.102436

[8] S. Teerakanok, T. Uehara, and A. Inomata, "Migrating to zero trust architecture: Reviews and challenges," *Security and Communication Networks*, vol. 2021, Article 9947347, 2021. doi:10.1155/2021/9947347

[9] M. L. Gambo and A. Almulhem, "Zero trust architecture: A systematic literature review," arXiv:2503.11659, 2025.

[10] S. Yuan and X. Wu, "Deep learning for insider threat detection: Review, challenges and opportunities," *Computers & Security*, vol. 104, 102221, 2021. doi:10.1016/j.cose.2021.102221

[11] D. C. Le and N. Zincir-Heywood, "Anomaly detection for insider threats using unsupervised ensembles," *IEEE Transactions on Network and Service Management*, vol. 18, no. 2, 2021.

[12] S. K. Basak, L. Neil, B. Reaves, and L. Williams, "What are the practices for secret management in software artifacts?," in *Proc. IEEE Secure Development Conference (SecDev)*, 2022. arXiv:2208.11280

[13] A. P. Singh and A. Sharma, "A systematic literature review on insider threats," arXiv:2212.05347, 2022.

[14] Verizon, "2025 Data Breach Investigations Report," 2025. https://www.verizon.com/business/resources/reports/dbir/

[15] NIST Joint Task Force, "Security and Privacy Controls for Information Systems and Organizations," NIST SP 800-53 Rev. 5, 2020. doi:10.6028/NIST.SP.800-53r5

[16] Gartner, "Magic Quadrant for Privileged Access Management," Oct. 2025 (M. Kelley et al.), as cited in CyberArk and BeyondTrust press releases, Oct. 2025: https://www.cyberark.com/press/cyberark-named-a-leader-in-the-2025-gartner-magic-quadrant-for-privileged-access-management/ ; Palo Alto Networks completion of CyberArk acquisition: https://www.techmonitor.ai/?p=419881

[17] J. Glasser and B. Lindauer, "Bridging the gap: A pragmatic approach to generating insider threat data," in *Proc. IEEE Security and Privacy Workshops*, 2013, pp. 98–104. doi:10.1109/SPW.2013.37

[18] V. C. Hu et al., "Guide to Attribute Based Access Control (ABAC) Definition and Considerations," NIST SP 800-162, 2014. doi:10.6028/NIST.SP.800-162

> **Before the pre-defense:** open each DOI once and confirm page numbers against the publisher record; add any sources your supervisor recommends; keep this list in sync with the thesis bibliography.

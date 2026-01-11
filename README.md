# M1 Middleware - Timetable Alerter (TPGO)

Ce projet est une application distribuée en **Golang** réalisée dans le cadre du cours de Middleware. Son objectif est de surveiller les emplois du temps universitaires (format iCal) et d'alerter les étudiants par email en cas de modification (changement de salle, nouvel horaire, suppression), en utilisant une architecture **microservices** et une communication **asynchrone** via **NATS JetStream**.

🔗 **Lien du dépôt GitHub :** [https://github.com/LovisCoding/M1_MIDDLEWARE_TPGO](https://github.com/LovisCoding/M1_MIDDLEWARE_TPGO)

---

## 🏗 Architecture du projet

Le projet est divisé en plusieurs microservices autonomes qui communiquent entre eux :

1. **Config API (`/config`)** :
   - API REST de gestion des configurations.
   - Permet de définir les agendas (ID des groupes) et les règles d'alertes (qui prévenir pour quel agenda).
   - Stockage : SQLite.

2. **Timetable API & Consumer (`/go-timetable`)** :
   - Gère les événements (cours) et détecte les changements.
   - **Consumer** : Écoute les événements bruts envoyés par le Scheduler, les compare avec sa base de données locale, et publie un message si une différence est détectée.
   - Stockage : SQLite.

3. **Scheduler (`/scheduler`)** :
   - Service autonome qui récupère périodiquement les fichiers iCal depuis les serveurs de l'UCA.
   - Parse les données et les envoie dans la file de messages (NATS).

4. **Alerter (`/alerter`)** :
   - **Consumer** : Écoute les messages de "modification validée" envoyés par le service Timetable.
   - Récupère les destinataires via l'API Config et envoie les emails de notification via une API de mail externe.

5. **Infrastructure (NATS)** :
   - Serveur de messagerie (Message Queue) assurant la communication asynchrone entre le Scheduler, le Timetable et l'Alerter.

6. **Frontend (`/frontend`)** :
   - Application web développée avec **SvelteKit** et **Vite**.
   - Interface utilisateur permettant de visualiser l'état du système et potentiellement configurer les alertes.

---

## 🚀 Prérequis

Avant de lancer le projet, assurez-vous d'avoir installé :

- **Go** (version 1.21 ou supérieure recommandée)
- **Node.js** (pour lancer le frontend)
- **Docker** (pour lancer le serveur NATS)
- **Make** (pour utiliser les commandes d'automatisation)

---

## 🛠 Installation et Lancement

### 1. Démarrer l'infrastructure (NATS)

Le projet nécessite un serveur NATS avec JetStream activé. Une commande `make` est configurée pour lancer un conteneur Docker automatiquement.

    make nats

*Cela va télécharger l'image docker `nats`, lancer le conteneur sur le port `4222`.*

### 2. Lancer les services

Vous avez plusieurs options pour lancer le projet selon vos besoins.

#### Option A : Tout lancer en même temps (Recommandé)
Cette commande compile et lance tous les services (`config`, `scheduler`, `alerter`, etc.) en parallèle. Les logs seront mélangés dans la console.

    make run-all

#### Option B : Lancer les services individuellement
Vous pouvez ouvrir plusieurs terminaux et lancer chaque brique séparément pour mieux voir les logs de chaque service :

- **Terminal 1 (Config) :**
    make config

- **Terminal 2 (Scheduler) :**
    make scheduler

- **Terminal 3 (Alerter) :**
    make alerter

- **Terminal 4 (Frontend) :**
    make frontend

### 3. Arrêter le projet

Pour stopper les services lancés avec `run-all`, faites simplement `Ctrl+C`.

---

## 📚 Documentation API (Swagger)

Chaque API REST dispose de sa propre documentation Swagger auto-générée.

Pour régénérer la documentation (si vous avez modifié le code), vous pouvez utiliser la commande suivante dans le dossier du service concerné (`config/` ou `go-timetable/`) :

    make swag

*Nécessite l'outil `swag` : `go install github.com/swaggo/swag/cmd/swag@latest`*

---

## 👥 Auteurs

Projet réalisé par Arthur LECOMTE et Flavien BOUHAMDANI.
#!/usr/bin/env python3
"""
Simple ABNF generator for taskspec grammar.
This script generates random valid taskspec entries based on the grammar.
"""

import random
import sys
import time
from datetime import datetime, timedelta

# Keywords
KEYWORDS = ["TODO", "FIXME", "BUG", "HACK", "NOTE", "INFO", "IDEA", "REFACTOR", "REMINDER"]

# Priority values
PRIORITY_VALUES = ["highest", "critical", "1", "high", "2", "medium", "normal", "3", "low", "4", "lowest", "5"]
PRIORITY_EMOJIS = ["🔺", "⏫", "🔼", "🔽", "⏬"]

# Status values
STATUS_VALUES = ["todo", "in-progress", "done", "cancelled", "blocked"]
STATUS_EMOJIS = ["⬜", "🚧", "✅", "❌", "🚫"]

# Sample names and tags
NAMES = ["alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi"]
TAGS = ["backend", "frontend", "urgent", "bug", "feature", "refactor", "test"]
PROJECTS = ["ProjectX", "ProjectY", "Alpha", "Beta", "Gamma"]

def random_date():
    """Generate a random date in the future"""
    days = random.randint(1, 365)
    future = datetime.now() + timedelta(days=days)
    return future.strftime("%Y-%m-%d")

def random_datetime():
    """Generate a random datetime in the future"""
    days = random.randint(1, 365)
    hours = random.randint(0, 23)
    minutes = random.randint(0, 59)
    seconds = random.randint(0, 59)
    future = datetime.now() + timedelta(days=days, hours=hours, minutes=minutes, seconds=seconds)
    return future.strftime("%Y-%m-%dT%H:%M:%SZ")

def random_description():
    """Generate a random description"""
    descriptions = [
        "Fix the bug",
        "Implement new feature",
        "Refactor code",
        "Update documentation",
        "Add tests",
        "Deploy to production",
        "Review pull request",
        "Optimize performance",
        "Fix security vulnerability",
        "Update dependencies",
        "Buy groceries",
        "Complete report",
        "Call the client",
        "Schedule meeting",
    ]
    return random.choice(descriptions)

def random_recurrence():
    """Generate a random recurrence pattern"""
    patterns = [
        "every day",
        "every week",
        "every month",
        "every 2 days",
        "every Monday",
        "daily",
        "weekly",
        "monthly",
    ]
    return random.choice(patterns)

def random_identifier():
    """Generate a random identifier"""
    prefixes = ["TASK", "BUG", "FEATURE", "ISSUE"]
    return f"{random.choice(prefixes)}-{random.randint(100, 9999)}"

def random_estimate():
    """Generate a random time estimate"""
    estimates = ["30m", "1h", "2h", "4h", "1d", "2d", "1w", "2w"]
    return random.choice(estimates)

def generate_metadata():
    """Generate random metadata fields"""
    metadata = []
    
    # Randomly add due date (30% chance)
    if random.random() < 0.3:
        if random.random() < 0.5:
            metadata.append(f"due:{random_date()}")
        else:
            metadata.append(f"📅 {random_date()}")
    
    # Randomly add scheduled date (20% chance)
    if random.random() < 0.2:
        if random.random() < 0.5:
            metadata.append(f"scheduled:{random_date()}")
        else:
            metadata.append(f"⏳ {random_date()}")
    
    # Randomly add start date (15% chance)
    if random.random() < 0.15:
        if random.random() < 0.5:
            metadata.append(f"start:{random_date()}")
        else:
            metadata.append(f"🛫 {random_date()}")
    
    # Randomly add priority (40% chance)
    if random.random() < 0.4:
        if random.random() < 0.3:
            metadata.append(random.choice(PRIORITY_EMOJIS))
        elif random.random() < 0.5:
            metadata.append(f"priority:{random.choice(PRIORITY_VALUES)}")
        else:
            metadata.append(f"p:{random.choice(PRIORITY_VALUES)}")
    
    # Randomly add assignees (30% chance)
    if random.random() < 0.3:
        num_assignees = random.randint(1, 2)
        for _ in range(num_assignees):
            if random.random() < 0.8:
                metadata.append(f"@{random.choice(NAMES)}")
            else:
                metadata.append(f"👤{random.choice(NAMES)}")
    
    # Randomly add tags (40% chance)
    if random.random() < 0.4:
        num_tags = random.randint(1, 3)
        for _ in range(num_tags):
            metadata.append(f"#{random.choice(TAGS)}")
    
    # Randomly add projects (20% chance)
    if random.random() < 0.2:
        metadata.append(f"+{random.choice(PROJECTS)}")
    
    # Randomly add status (25% chance)
    if random.random() < 0.25:
        if random.random() < 0.3:
            metadata.append(random.choice(STATUS_EMOJIS))
        else:
            metadata.append(f"status:{random.choice(STATUS_VALUES)}")
    
    # Randomly add recurrence (15% chance)
    if random.random() < 0.15:
        if random.random() < 0.5:
            metadata.append(f"repeat:{random_recurrence()}")
        else:
            metadata.append(f"🔁 {random_recurrence()}")
    
    # Randomly add identifier (25% chance)
    if random.random() < 0.25:
        if random.random() < 0.5:
            metadata.append(f"id:{random_identifier()}")
        else:
            metadata.append(f"🆔 {random_identifier()}")
    
    # Randomly add estimate (20% chance)
    if random.random() < 0.2:
        if random.random() < 0.5:
            metadata.append(f"estimate:{random_estimate()}")
        else:
            metadata.append(f"⏱️ {random_estimate()}")
    
    return metadata

def generate_generic_todo():
    """Generate a generic TODO entry"""
    keyword = random.choice(KEYWORDS)
    description = random_description()
    metadata = generate_metadata()
    
    result = f"{keyword}: {description}"
    if metadata:
        result += " " + " ".join(metadata)
    return result

def generate_markdown_todo():
    """Generate a Markdown task list entry"""
    # Use lowercase x to match ABNF grammar specification
    checked = random.choice(["[ ]", "[x]"])
    description = random_description()
    metadata = generate_metadata()
    
    result = f"- {checked} {description}"
    if metadata:
        result += " " + " ".join(metadata)
    return result

def generate_sample():
    """Generate a random taskspec sample"""
    # 70% generic, 30% markdown
    if random.random() < 0.7:
        return generate_generic_todo()
    else:
        return generate_markdown_todo()

def main():
    """Main function to generate samples"""
    if len(sys.argv) > 1:
        try:
            duration_seconds = int(sys.argv[1])
        except ValueError:
            print(f"Usage: {sys.argv[0]} [duration_in_seconds]", file=sys.stderr)
            sys.exit(1)
    else:
        duration_seconds = 600  # Default to 10 minutes
    
    start_time = time.time()
    count = 0
    
    while time.time() - start_time < duration_seconds:
        sample = generate_sample()
        print(sample)
        count += 1
        
        # Small delay to avoid overwhelming the output
        if count % 1000 == 0:
            sys.stderr.write(f"Generated {count} samples...\n")
            sys.stderr.flush()
    
    sys.stderr.write(f"Total samples generated: {count}\n")

if __name__ == "__main__":
    main()

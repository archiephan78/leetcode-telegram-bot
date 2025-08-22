#!/usr/bin/env python3
"""
Script to remove duplicate LeetCode problems from problem.yaml file.
Keeps the first occurrence of each problem and removes subsequent duplicates.
"""

from typing import Dict, List, Set, Any
import sys


def deduplicate_problems(input_file: str, output_file: str) -> None:
    """
    Remove duplicate problems from the YAML file.
    
    Args:
        input_file: Path to the input YAML file
        output_file: Path to the output deduplicated YAML file
    """
    # Read the original file content
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Parse the content manually since it's not standard YAML
    categories = {}
    current_category = None
    seen_problems: Set[str] = set()
    
    lines = content.strip().split('\n')
    
    for line in lines:
        line = line.rstrip()
        
        # Skip empty lines
        if not line:
            continue
            
        # Check if this is a category line (doesn't start with - or space)
        if not line.startswith((' ', '-', '\t')):
            current_category = line.rstrip(':')
            categories[current_category] = []
        elif line.startswith('- title:'):
            # Extract the title
            title = line.split('- title: ', 1)[1].strip()
            
            # Check if we've seen this problem before
            if title not in seen_problems:
                seen_problems.add(title)
                categories[current_category].append({'title': title, 'pending_url': True})
            else:
                print(f"Removing duplicate: {title} from category: {current_category}")
        elif line.strip().startswith('url:') and current_category and categories[current_category]:
            # Add URL to the last problem if it doesn't have one yet
            last_problem = categories[current_category][-1]
            if last_problem.get('pending_url'):
                url = line.strip().split('url: ', 1)[1]
                last_problem['url'] = url
                del last_problem['pending_url']
        elif line.strip().startswith('- title:') and current_category:
            # Handle nested categories (like "Other Advanced Patterns")
            nested_title = line.strip().split('- title: ', 1)[1]
            if nested_title not in seen_problems:
                seen_problems.add(nested_title)
                categories[current_category].append({'title': nested_title, 'pending_url': True})
            else:
                print(f"Removing duplicate: {nested_title} from category: {current_category}")
    
    # Write the deduplicated content
    with open(output_file, 'w', encoding='utf-8') as f:
        for category, problems in categories.items():
            if not problems:  # Skip empty categories
                continue
                
            f.write(f"{category}:\n")
            
            for problem in problems:
                if 'pending_url' in problem:
                    continue  # Skip problems without URLs
                    
                f.write(f"- title: {problem['title']}\n")
                f.write(f"  url: {problem['url']}\n")
    
    print(f"Deduplication complete. Original: {len(content.split('- title:')) - 1} problems")
    print(f"After deduplication: {len(seen_problems)} unique problems")


if __name__ == "__main__":
    input_file = "problem.yaml"
    output_file = "problem_deduplicated.yaml"
    
    try:
        deduplicate_problems(input_file, output_file)
        print(f"Deduplicated file saved as: {output_file}")
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1) 
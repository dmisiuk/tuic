#!/bin/bash

echo "Getting status..."
echo ""
echo ""

epic_name="$1"

if [ -z "$epic_name" ]; then
  echo "❌ Please specify an epic name"
  echo "Usage: /pm:epic-status <epic-name>"
  echo ""
  echo "Available epics:"
  for dir in .claude/epics/*/; do
    [ -d "$dir" ] && echo "  • $(basename "$dir")"
  done
  exit 1
else
  # Show status for specific epic
  epic_dir=".claude/epics/$epic_name"
  epic_file="$epic_dir/epic.md"

  if [ ! -f "$epic_file" ]; then
    echo "❌ Epic not found: $epic_name"
    echo ""
    echo "Available epics:"
    for dir in .claude/epics/*/; do
      [ -d "$dir" ] && echo "  • $(basename "$dir")"
    done
    exit 1
  fi

  echo "📚 Epic Status: $epic_name"
  echo "================================"
  echo ""

  # Extract metadata
  status=$(grep "^status:" "$epic_file" | head -1 | sed 's/^status: *//')
  progress=$(grep "^progress:" "$epic_file" | head -1 | sed 's/^progress: *//')
  github=$(grep "^github:" "$epic_file" | head -1 | sed 's/^github: *//')

  # Count tasks
  total=0
  open=0
  closed=0
  blocked=0

  # Use find to safely iterate over task files (excluding analysis files)
  for task_file in "$epic_dir"/[0-9]*.md; do
    [ -f "$task_file" ] || continue
    # Skip analysis files
    [[ "$task_file" =~ -analysis\.md$ ]] && continue
    ((total++))

    task_status=$(grep "^status:" "$task_file" | head -1 | sed 's/^status: *//')
    deps=$(grep "^depends_on:" "$task_file" | head -1 | sed 's/^depends_on: *\[//' | sed 's/\]//' | sed 's/depends_on: *//')

    if [ "$task_status" = "closed" ] || [ "$task_status" = "completed" ]; then
      ((closed++))
    elif [ -n "$deps" ] && [ "$deps" != "[]" ] && [ "$deps" != "" ]; then
      # Check if dependencies are actually blocking
      deps_met=true
      for dep in $deps; do
        # Remove commas and check if dependency task exists and is closed
        dep_clean=$(echo "$dep" | sed 's/,//')
        if [ "$dep_clean" = "001" ] && ! grep -q "status: closed" "$epic_dir/5.md"; then deps_met=false; fi
        if [ "$dep_clean" = "002" ] && ! grep -q "status: closed" "$epic_dir/8.md"; then deps_met=false; fi
        if [ "$dep_clean" = "003" ] && ! grep -q "status: closed" "$epic_dir/9.md"; then deps_met=false; fi
        if [ "$dep_clean" = "004" ] && ! grep -q "status: closed" "$epic_dir/10.md"; then deps_met=false; fi
        if [ "$dep_clean" = "005" ] && ! grep -q "status: closed" "$epic_dir/11.md"; then deps_met=false; fi
        if [ "$dep_clean" = "006" ] && ! grep -q "status: closed" "$epic_dir/12.md"; then deps_met=false; fi
        if [ "$dep_clean" = "007" ] && ! grep -q "status: closed" "$epic_dir/13.md"; then deps_met=false; fi
        if [ "$dep_clean" = "008" ] && ! grep -q "status: closed" "$epic_dir/6.md"; then deps_met=false; fi
      done
      if [ "$deps_met" = false ]; then
        ((blocked++))
      else
        ((open++))
      fi
    else
      ((open++))
    fi
  done

  # Display progress bar
  if [ $total -gt 0 ]; then
    percent=$((closed * 100 / total))
    filled=$((percent * 20 / 100))
    empty=$((20 - filled))

    echo -n "Progress: ["
    [ $filled -gt 0 ] && printf '%0.s█' $(seq 1 $filled)
    [ $empty -gt 0 ] && printf '%0.s░' $(seq 1 $empty)
    echo "] $percent%"
  else
    echo "Progress: No tasks created"
  fi

  echo ""
  echo "📊 Breakdown:"
  echo "  Total tasks: $total"
  echo "  ✅ Completed: $closed"
  echo "  🔄 Available: $open"
  echo "  ⏸️ Blocked: $blocked"

  [ -n "$github" ] && echo ""
  [ -n "$github" ] && echo "🔗 GitHub: $github"
fi

exit 0

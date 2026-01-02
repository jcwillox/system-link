
??? example "Dashboard configuration"

    I use [Mushroom](https://github.com/piitaya/lovelace-mushroom) cards here, but you could easily recreate this with standard tile cards.

    ```yaml
    type: grid
    cards:
      - type: heading
        heading: CommBank
        heading_style: title
        icon: mdi:bank
      - type: custom:mushroom-chips-card
        chips:
          - type: entity
            entity: switch.pyscripts_test_script
            icon_color: primary
            icon: mdi:script-text-play
          - type: entity
            entity: button.pyscripts_test_script_run
            tap_action:
              action: more-info
            hold_action:
              action: none
            double_tap_action:
              action: none
            content_info: none
            name: Run
            icon_color: red
            icon: ""
      - type: tile
        entity: sensor.pyscripts_test_script_duration
        name: Duration
        color: primary
        hide_state: false
        vertical: false
        features_position: bottom
      - type: tile
        entity: binary_sensor.pyscripts_test_script_successful
        name: Status
        vertical: false
        features_position: bottom
      - type: tile
        entity: sensor.pyscripts_test_script_last_run
        name: Last Run
        color: primary
        vertical: false
        features_position: bottom
      - type: tile
        entity: sensor.pyscripts_test_script_next_run
        name: Next Run
        color: primary
        vertical: false
        features_position: bottom
      - type: markdown
        content: >
          {% set entity_id = "binary_sensor.pyscripts_test_script_successful"
          %}


          {% set stderr = state_attr(entity_id, "stderr") %}

          {% set stdout = state_attr(entity_id, "stdout") %}


          {% if stdout %}

          #### Output

          ```

          {{ stdout }}

          ```

          {% endif %}


          {% if stderr %}

          #### Error

          ```

          {{ stderr }}

          ```

          {% endif %}


          #### Exit Code: {{ state_attr(entity_id, "code") or 0 }}
        visibility:
          - condition: state
            entity: binary_sensor.pyscripts_test_script_successful
            state: "on"
    ```

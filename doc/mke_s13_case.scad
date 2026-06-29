// =====================================================
// MKE-S13 Capacitive Soil Moisture Sensor — Case
// Parametric OpenSCAD  |  v5.1 (Drop-In U-Slot & Pillars)
// Units: mm
// =====================================================

include <mke_s13_config.scad>

$fn = 48;

// =====================================================
// 2D PROFILES
// =====================================================
module pcb_2d() {
    hull() {
        translate([0,        pcb_w/2])         circle(r=0.1);
        translate([chev_l,   0])                circle(r=0.1);
        translate([chev_l,   pcb_w])           circle(r=0.1);
        translate([pcb_l - corner_r, corner_r])           circle(r=corner_r);
        translate([pcb_l - corner_r, pcb_w - corner_r])   circle(r=corner_r);
    }
}

module pcb_box_section_2d() {
    intersection() {
        pcb_2d();
        translate([safe_line_x - 1, -5])
            square([box_l + 6, pcb_w + 10]);
    }
}

module shell_2d(w) {
    offset(r = w + pcb_gap) pcb_box_section_2d();
}

// =====================================================
// BOTTOM SHELL (Clean exterior, Drop-in U-Slot)
// =====================================================
module bottom_shell() {
    difference() {
        // Main Outer body box
        linear_extrude(outer_h) shell_2d(wall);

        // Interior hollow
        translate([0, 0, floor_t])
            linear_extrude(outer_h)
                offset(r = pcb_gap) pcb_box_section_2d();

        // --- DROP-IN U-SLOT FOR PCB ---
        // Open to the top rim to allow Z-axis insertion
        translate([
            safe_line_x - 1 - pcb_gap - wall - 1.0,
            pcb_w/2 - (pcb_w/2 + slot_gap),
            floor_t + pcb_t + 0.2 - slot_gap
        ])
            cube([
                wall + 2.0,
                pcb_w + 2*slot_gap,
                outer_h
            ]);

        // Pin alignment holes through the floor
        for (y_off = [-hole_sp/2, hole_sp/2]) {
            translate([hole_x, hole_cy + y_off, -0.01])
                cylinder(d = hole_d, h = floor_t + 1);
        }

        // Alignment nub dimples on rim
        translate([nub_x, -wall*0.5 - pcb_gap, outer_h - nub_h]) cylinder(d = nub_d + nub_clearance, h = nub_h + 0.1);
        translate([nub_x, pcb_w + wall*0.5 + pcb_gap, outer_h - nub_h]) cylinder(d = nub_d + nub_clearance, h = nub_h + 0.1);
    }

    // PCB mounting standoffs (Pins from the lid will seat inside here)
    translate([0, 0, floor_t]) {
        for (y_off = [-hole_sp/2, hole_sp/2]) {
            translate([hole_x, hole_cy + y_off, 0])
            difference() {
                cylinder(d = pcb_boss_d, h = pcb_t + 0.2);
                translate([0, 0, -0.01])
                    cylinder(d = hole_d, h = pcb_t + 1);
            }
        }
    }
}

// =====================================================
// LID (Locking pillars, Rotated Slot, Closure Baffle)
// =====================================================
module lid() {
    difference() {
        union() {
            // Main Lid plate
            linear_extrude(lid_t) shell_2d(wall);

            // --- CLOSURE BAFFLE ---
            // Drops down into the bottom shell's U-slot to seal the case
            baffle_h = inner_h - (2 * pcb_t + 0.2 + slot_gap);
            baffle_w = pcb_w + 2*slot_gap - baffle_clearance; // baffle_clearance mm tolerance total for smooth fit

            translate([
                safe_line_x - 1 - pcb_gap - wall,
                pcb_w/2 - baffle_w/2,
                -baffle_h
            ])
                cube([wall, baffle_w, baffle_h]);
        }

        // ROTATED CONNECTOR SLOT (Wide along Y axis to fit the plug)
        translate([
            pcb_l - conn_d - cable_clear,
            pcb_w/2 - (conn_l / 2) - cable_clear,
            -0.01
        ])
            cube([
                conn_d + 2*cable_clear,
                conn_l + 2*cable_clear,
                lid_t + 0.1
            ]);
    }

    // Alignment registration nubs
    translate([nub_x, -wall*0.5 - pcb_gap, -nub_h]) cylinder(d = nub_d, h = nub_h);
    translate([nub_x, pcb_w + wall*0.5 + pcb_gap, -nub_h]) cylinder(d = nub_d, h = nub_h);

    // --- INTEGRATED LOCKING PILLARS ---
    // Protrudes downwards from the underside of the lid (Z=0 local)
    for (y_off = [-hole_sp/2, hole_sp/2]) {
        translate([hole_x, hole_cy + y_off, 0]) {

            // 1. Upper Wide Shoulder: Clamps onto the top face of the PCB substrate
            shoulder_h = inner_h - (floor_t + 2 * pcb_t + 0.2);
                    translate([0, 0, -shoulder_h])
                        cylinder(d = pcb_boss_d, h = shoulder_h);

            // 2. Interlocking Pin: Passes through PCB and bottom standoffs to the outside floor
            translate([0, 0, -inner_h])
                        cylinder(d = lock_pin_d, h = 2 * pcb_t + 0.2 + floor_t);
        }
    }
}

// =====================================================
// PRINT-READY ORIENTATIONS
// =====================================================
// The lid() module above is authored in "assembly logic":
// the flat plate sits at the TOP (z=0..lid_t) and the
// closure baffle + locking pillars hang DOWNWARD into
// negative Z, matching how the lid sits when placed onto
// the bottom shell. That orientation is NOT printable as-is
// (parts fall below the z=0 bed plane), and printing it
// pillars-down would require support material under every
// pillar tip.
//
// lid_print_ready() mirrors the part so the plate becomes
// the base sitting flat on the bed, and the pillars become
// self-supporting upward-growing towers — no supports needed.
// Print the plate first; the pillars build up naturally on
// top of it, same logic as the bottom shell's standoffs.
module lid_print_ready() {
    translate([0, 0, lid_t])
        mirror([0, 0, 1])
            lid();
}

// =====================================================
// RENDER
// =====================================================
color("SteelBlue",  0.85) bottom_shell();
color("LightBlue",  0.70) translate([0, pcb_w + 15, 0]) lid_print_ready();
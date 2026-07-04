// =====================================================
// MKE-S13 Case + PCB — Combined Fit-Check Assembly
// Parametric OpenSCAD  | v8.2 (Fixed CSG Wall Glitch in Terminal Holes)
// Units: mm
// =====================================================
include <BOSL2/std.scad>
include <BOSL2/screws.scad>
use <mke_s13_case_b.scad>
include <mke_s13_config_b.scad>

// --- VIEW CONFIGURATION ---
// 1 = Assembled (Closed with transparency)
// 2 = Exploded (Separated vertically)
// 3 = Cutaway (Sliced along Y-midplane to verify clearances)
view_mode = 3;

$fn = 64;

// =====================================================
// PCB SUB-COMPONENTS
// =====================================================
module pcb_outline_2d() {
    hull() {
        translate([0, pcb_w/2]) circle(r=0.1);
        translate([chev_l, 0]) circle(r=0.1);
        translate([chev_l, pcb_w]) circle(r=0.1);
        translate([pcb_l - corner_r, corner_r]) circle(r=corner_r);
        translate([pcb_l - corner_r, pcb_w - corner_r]) circle(r=corner_r);
    }
}

module pcb_board() {
    color("ForestGreen", 0.95) {
        difference() {
            linear_extrude(pcb_t) pcb_outline_2d();
            for(i = [-1, 0, 1]) {
                translate([pcb_l - pin_offset_x, pcb_w/2 + (i * conn_pitch), -0.5]) {
                    cylinder(d=1.0, h=pcb_t + 1);
                }
            }
            for (y_off = [-hole_sp/2, hole_sp/2]) {
                translate([hole_x, hole_cy + y_off, -0.5]) {
                    cylinder(d=hole_d, h=pcb_t + 1);
                }
            }
        }
    }
}

module connector_male(z_extra = 0) {
    color("White", 0.95) {
        translate([pcb_l - conn_male_d, pcb_w/2 - conn_male_w/2, pcb_t + z_extra]) {
            difference() {
                union() {
                    cube([conn_male_d, conn_male_w, conn_male_h]);
                    translate([conn_male_d, (conn_male_w/2) - 1.8, conn_male_h - 2.0]) {
                        hull() {
                            cube([0.6, 3.6, 0.01]);
                            translate([0, 0, 1.99]) {
                                cube([0.01, 3.6, 0.01]);
                            }
                        }
                    }
                }
                translate([0.8, 0.8, 1.5]) {
                    cube([conn_male_d - 1.6, conn_male_w - 1.6, conn_male_h]);
                }
                // Fix: coincident-face CSG glitch (slot's near face landed on
                // the exact same plane as the cavity cutout's far face) --
                // overlap slightly into the cavity so the two cuts genuinely
                // intersect instead of just touching.
                for(i = [-1, 1]) {
                    translate([conn_male_d - 0.8 - overlap_eps, (conn_male_w/2) + (i * conn_pitch) - 0.7, 1.5]) {
                        cube([1.0 + overlap_eps, 1.4, conn_male_h]);
                    }
                }
            }
        }
    }
    color("Silver") {
        pin_tip_h = 0.4;
        for(i = [-1, 0, 1]) {
            translate([pcb_l - pin_offset_x, pcb_w/2 + (i * conn_pitch), -pin_protrusion + z_extra]) {
                cylinder(d1=0, d2=0.64, h=pin_tip_h, $fn=16);
            }
            translate([pcb_l - pin_offset_x, pcb_w/2 + (i * conn_pitch), -pin_protrusion + pin_tip_h + z_extra]) {
                cylinder(d=0.64, h=(conn_male_h - 1.5) + pcb_t + pin_protrusion - pin_tip_h, $fn=16);
            }
        }
    }
}

module connector_female(z_extra = 0) {
    translate([0, 0, z_extra]) {
        color("Gainsboro", 0.98) {
            difference() {
                union() {
                    translate([pcb_l - conn_male_d + 0.85, pcb_w/2 - (conn_male_w - 1.7)/2, pcb_t + 1.5]) {
                        cube([conn_male_d - 1.7, conn_male_w - 1.7, 6.2]);
                        translate([-0.2, -0.4, 5.2]) {
                            cube([conn_male_d - 1.3, conn_male_w - 0.9, 1.2]);
                        }
                        for(i = [-1, 1]) {
                            translate([conn_male_d - 1.7, (conn_male_w - 1.7)/2 + (i * conn_pitch) - 0.6, 0]) {
                                cube([0.85, 1.2, 5.2]);
                            }
                        }
                        translate([conn_male_d - 1.7, (conn_male_w - 1.7)/2 - conn_pitch - 0.6, 5.2]) {
                            cube([0.85, (2 * conn_pitch) + 1.2, 1.2]);
                        }
                    }
                }

                for(i = [-1, 0, 1]) {
                    // A. Large Rectangle (Right side)
                    translate([
                        (pcb_l - 1.35) - 2.5,
                        (pcb_w / 2) + (i * conn_pitch) - 1.0,
                        pcb_t + 1.5 + 3.0
                    ]) {
                        cube([2.5, 2.0, 4.0]);
                    }

                    // B. Small Rectangle (Left side)
                    // 0.5mm nominal length (X) + terminal_hole_overlap into the
                    // large rectangle to prevent the CGAL coincident-face glitch
                    translate([
                        (pcb_l - 1.35) - 2.5 - 0.5,
                        (pcb_w / 2) + (i * conn_pitch) - 0.5,
                        pcb_t + 1.5 + 3.0
                    ]) {
                        cube([0.5 + terminal_hole_overlap, 1.0, 4.0]);
                    }
                }
            }
        }

        wire_colors = ["Red", "White", "Black"];
        for(i = [-1, 0, 1]) {
            color(wire_colors[i+1]) {
                translate([
                    pcb_l - pin_offset_x,
                    pcb_w/2 + (i * conn_pitch),
                    pcb_t + 1.5 + 3.0
                ]) {
                    cylinder(d=1.1, h=20.0, $fn=16);
                }
            }
        }
    }
}

module safe_line_marker() {
    color("Black") {
        translate([safe_line_x, 0, pcb_t]) {
            hull() {
                translate([0, 3.25, 0])
                    cylinder(d=1.5, h=0.05, $fn=32);
                translate([0, pcb_w - 3.25, 0])
                    cylinder(d=1.5, h=0.05, $fn=32);
            }
        }
    }
}

module red_zone_markers() {
    color("Red") {
        translate([pcb_l - red_line_near_edge - 0.75, 1, pcb_t]) {
            cube([1.5, pcb_w - 2, 0.05]);
        }
        translate([pcb_l - red_line_far_edge - 0.75, 1, pcb_t]) {
            cube([1.5, pcb_w - 2, 0.05]);
        }
    }
}

module connector_footprint_outline() {
    ox = pcb_l - conn_male_d;
    oy = pcb_w / 2 - conn_male_w / 2;
    silk_t = nozzle_d;
    silk_z = 0.05;

    color("White") translate([ox, oy, pcb_t]) {
        translate([0, 0, 0]) cube([conn_male_d, silk_t, silk_z]);
        translate([0, conn_male_w - silk_t, 0]) cube([conn_male_d, silk_t, silk_z]);
        translate([0, silk_t, 0]) cube([silk_t, conn_male_w - 2 * silk_t, silk_z]);
        translate([conn_male_d - silk_t, silk_t, 0]) cube([silk_t, conn_male_w - 2 * silk_t, silk_z]);
    }
}

module case_screws(z_lift = 0) {
    z_lid_top = outer_h + lid_t + z_lift;
    color("Silver") {
        for (p = corner_positions)
            translate([p[0], p[1], z_lid_top])
                screw("M3", length=12, head="flat", anchor=TOP);
    }
}

module pcb_assembly(male_z_extra = 0, female_z_extra = 0) {
    pcb_board();
    connector_male(z_extra = male_z_extra);
    connector_female(z_extra = female_z_extra);
    safe_line_marker();
    red_zone_markers();
    connector_footprint_outline();
}

module full_system() {
    color("SteelBlue", 0.65) bottom_shell();
    translate([0, 0, z_pcb_seat]) pcb_assembly();
    translate([0, 0, outer_h]) color("LightBlue", 0.50) lid();
    case_screws(z_lift = 0);
}

module exploded_system() {
    gap = 25;
    color("SteelBlue", 0.85) bottom_shell();
    translate([0, 0, z_pcb_seat + gap]) {
        pcb_board();
        safe_line_marker();
        red_zone_markers();
        connector_footprint_outline();
    }
    translate([0, 0, z_pcb_seat + gap]) {
        connector_male(z_extra = gap);
    }
    translate([0, 0, z_pcb_seat + gap]) {
        connector_female(z_extra = gap + conn_male_h + gap);
    }
    translate([0, 0, outer_h + gap * 4]) color("LightBlue", 0.75) lid();
    case_screws(z_lift = 25 * 5);
}

if (view_mode == 1) {
    full_system();
} else if (view_mode == 2) {
    exploded_system();
} else if (view_mode == 3) {
    difference() {
        full_system();
        translate([-10, pcb_w/2, -5]) {
            cube([pcb_l + 30, pcb_w + 20, outer_h + 30]);
        }
    }
}